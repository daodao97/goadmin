package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"log"
	"os"
)

var configFile = flag.String("c", "application.toml", "application config file")

// Conf 服务配置
type Conf struct {
	// 脚手架必须配置
	scaffold.Conf
    // 自定义配置
    Redis *redis.Options
    
}

func NewConf() (c *Conf, cf func(), err error) {
	conf, err := os.ReadFile(*configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return
	}
	_, err = toml.Decode(string(conf), &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return
	}
	c1, err := initDb(c)
	if err != nil {
		return nil, func() { c1() }, err
	}
	return c, func() {
		c1()
	}, nil

	return c, func() {}, nil
}

func initDb(c *Conf) (func(), error) {
	conns := make(map[string]*db.Config)
	for _, v := range c.Database {
		sqlConf := new(db.Config)
		err := copier.Copy(sqlConf, v)
		if err != nil {
			return func() {}, err
		}
		conns[v.Name] = sqlConf
	}
	err := db.Init(conns)
	if err != nil {
		return nil, err
	}
	return func() {
		db.Close()
	}, nil
}
