package conf

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "", "config path, example: -conf /config.yaml")
}

// 解析配置文件
func Parse(c interface{}) error {
	flag.Parse()
	if confPath == "" {
		return errors.New("load config file path failed")
	}
	viper.SetConfigFile(confPath)
	viper.SetConfigType("yaml")
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
