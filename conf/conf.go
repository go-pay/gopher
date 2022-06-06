package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	_Json = iota + 1
	_Yaml
	_Toml
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "", "config path, example: -conf /config.yaml")
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseYaml(confPtr interface{}) error {
	return parse(_Yaml, confPtr)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseToml(confPtr interface{}) error {
	return parse(_Toml, confPtr)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseJson(confPtr interface{}) error {
	return parse(_Json, confPtr)
}

func parse(cType int, confPtr interface{}) error {
	if confPtr == nil {
		return errors.New("c struct ptr can not be nil")
	}

	beanValue := reflect.ValueOf(confPtr)
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
	fileBs, err := ioutil.ReadFile(confPath)
	if err != nil {
		return fmt.Errorf("read config file error: %w", err)
	}
	switch cType {
	case _Json:
		if err = json.Unmarshal(fileBs, confPtr); err != nil {
			return fmt.Errorf("parse config file [%s] error: %w", string(fileBs), err)
		}
	case _Yaml:
		if err = yaml.Unmarshal(fileBs, confPtr); err != nil {
			return fmt.Errorf("parse config file [%s] error: %w", string(fileBs), err)
		}
	case _Toml:
		if _, err = toml.Decode(string(fileBs), confPtr); err != nil {
			return fmt.Errorf("parse config file [%s] error: %w", string(fileBs), err)
		}
	default:
		return errors.New("config file only support: yaml、json、toml")
	}
	return nil
}
