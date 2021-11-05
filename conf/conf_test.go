package conf

import (
	"flag"
	"testing"

	"github.com/go-pay/gopher/orm"
	"github.com/go-pay/gopher/web"
	"github.com/go-pay/gopher/xlog"
)

type Config struct {
	Name   string           `json:"name" yaml:"name" toml:"name"`
	Number int              `json:"number" yaml:"number" toml:"number"`
	Web    *web.Config      `json:"web" yaml:"web" toml:"web"`
	MySQL  *orm.MySQLConfig `json:"mysql" yaml:"mysql" toml:"mysql"`
	Redis  *orm.RedisConfig `json:"redis" yaml:"redis" toml:"redis"`
}

func TestParseYaml(t *testing.T) {
	c := &Config{}
	flag.Set("conf", "config.yaml")
	if err := ParseYaml(c); err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug(c.Name)
	xlog.Debug(c.Number)
	xlog.Debug(c.Web)
	xlog.Debug(c.Web.Limiter)
	xlog.Debug(c.MySQL)
	xlog.Debug(c.Redis)
}

func TestParseJson(t *testing.T) {
	c := &Config{}
	flag.Set("conf", "config.json")
	if err := ParseJson(c); err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug(c.Name)
	xlog.Debug(c.Number)
	xlog.Debug(c.Web)
	xlog.Debug(c.Web.Limiter)
	xlog.Debug(c.MySQL)
	xlog.Debug(c.Redis)
}

func TestParseToml(t *testing.T) {
	c := &Config{}
	flag.Set("conf", "config.toml")
	if err := ParseToml(c); err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug(c.Name)
	xlog.Debug(c.Number)
	xlog.Debug(c.Web)
	xlog.Debug(c.Web.Limiter)
	xlog.Debug(c.MySQL)
	xlog.Debug(c.Redis)
}
