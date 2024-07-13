package conf

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
)

var configFile = flag.String("c", "application.toml", "application config file")

var conf *Conf

// Conf 服务配置
type Conf struct {
	scaffold.Conf
	Redis *redis.Options
}

var scaffoldConfig *scaffold.Conf

func Get() *Conf {
	return conf
}

func GetScaffoldConf() *scaffold.Conf {
	return scaffoldConfig
}

func Startup() error {
	_conf, err := os.ReadFile(*configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	_, err = toml.Decode(string(_conf), &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	conns := make(map[string]*db.Config)
	for _, v := range conf.Database {
		sqlConf := new(db.Config)
		err := copier.Copy(sqlConf, v)
		if err != nil {
			return err
		}
		conns[v.Name] = sqlConf
	}
	err = db.Init(conns)
	if err != nil {
		return err
	}

	sc, err := scaffoldConf(conf)
	if err != nil {
		return err
	}

	scaffoldConfig = sc

	return nil
}

func scaffoldConf(c *Conf) (*scaffold.Conf, error) {
	cf := &scaffold.Conf{}
	err := copier.Copy(cf, c)
	if err != nil {
		return nil, err
	}
	err = cf.Validate()
	if err != nil {
		return nil, err
	}

	return cf, nil
}
