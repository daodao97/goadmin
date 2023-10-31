package util

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

func Md5(str string) string {
	data := []byte(str)
	md5Ctx := md5.New()
	_, _ = md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

func Length(v interface{}) int {
	switch v := v.(type) {
	case string:
		return len(v)
	case []interface{}:
		return len(v)
	case []int:
		return len(v)
	case []int64:
		return len(v)
	case []float64:
		return len(v)
	case []float32:
		return len(v)
	}
	return 0
}

func RecursiveDir(pathname string, f func(filePath string)) {
	rd, err := ioutil.ReadDir(pathname)
	if err == nil {
		for _, fi := range rd {
			if fi.IsDir() {
				RecursiveDir(pathname+fi.Name()+"\\", f)
			} else {
				f(pathname + "/" + fi.Name())
			}
		}
	}
}

func DirectoryExists(path string) bool {
	// Use os.Stat to get file info for the given path
	// If the directory exists, os.Stat will return nil (no error)
	// If it doesn't exist, os.IsNotExist will return true
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // Directory doesn't exist
		}
		// Some other error occurred (e.g., permission issue)
		// You may choose to handle it differently based on your use case
		return false
	}
	return true // Directory exists
}
