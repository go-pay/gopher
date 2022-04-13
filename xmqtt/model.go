package xmqtt

import (
	"errors"
	"time"
)

const (
	QosAtMostOne  QosType = 0
	QosAtLeastOne QosType = 1
	QosOnlyOne    QosType = 2
)

var (
	ErrNotExists   = errors.New("mqtt not exists")
	ErrLostConnect = errors.New("mqtt connection lost")
)

type QosType byte

type Config struct {
	Broker       string        `json:"broker" yaml:"broker" toml:"broker"`
	TcpPort      int           `json:"tcp_port" yaml:"tcp_port" toml:"tcp_port"`
	ClientId     string        `json:"client_id" yaml:"client_id" toml:"client_id"`
	Uname        string        `json:"uname" yaml:"uname" toml:"uname"`
	Password     string        `json:"password" yaml:"password" toml:"password"`
	KeepAlive    time.Duration `json:"keep_alive" yaml:"keep_alive" toml:"keep_alive"` // 单位秒
	CleanSession bool          `json:"clean_session" yaml:"clean_session" toml:"clean_session"`
}
