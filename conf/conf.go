package conf

import (
	"flag"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	json = iota + 1
	yaml
	toml
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "", "config path, example: -conf /config.yaml")
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseYaml(c interface{}) error {
	return parse(yaml, c)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseToml(c interface{}) error {
	return parse(toml, c)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseJson(c interface{}) error {
	return parse(json, c)
}

func parse(cType int, c interface{}) error {
	if c == nil {
		return errors.New("c struct ptr can not be nil")
	}
	var cFileType string
	beanValue := reflect.ValueOf(c)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("c must be ptr")
	}
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("c must be struct ptr")
	}
	flag.Parse()
	if confPath == "" {
		return errors.New("load config file path failed, add arguments -conf ")
	}
	viper.SetConfigFile(confPath)
	switch cType {
	case json:
		cFileType = "json"
	case yaml:
		cFileType = "yaml"
	case toml:
		cFileType = "toml"
	default:
		return errors.New("config file only support: yaml、json、toml")
	}
	viper.SetConfigType(cFileType)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}
