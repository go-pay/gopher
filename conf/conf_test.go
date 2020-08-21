package conf

import (
	"flag"
	"testing"

	"github.com/iGoogle-ink/gotil/orm"
	"github.com/iGoogle-ink/gotil/web"
	"github.com/iGoogle-ink/gotil/xlog"
)

type Config struct {
	Name   string
	Number int
	Web    *web.Config
	MySQL  *orm.MySQLConfig
	Redis  *orm.RedisConfig
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
	xlog.Debug(c.Web.Limit)
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
	xlog.Debug(c.Web.Limit)
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
	xlog.Debug(c.Web.Limit)
	xlog.Debug(c.MySQL)
	xlog.Debug(c.Redis)
}
