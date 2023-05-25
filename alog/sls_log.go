package alog

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/go-pay/gopher/util"
)

type Client struct {
	Config   *Config
	Producer *producer.Producer
}

// 阿里云日志配置结构体
type Config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Project   string
	LogStore  string
	HostName  string
}

var (
	ErrInvalidParam = errors.New("config Missing parameter")
	ErrJsonStr      = errors.New("error json string")
)

// 初始化日志
func New(config *Config) (client *Client, err error) {
	hostname, _ := os.Hostname()
	config.HostName = hostname
	err = checkConfig(config)
	if err != nil {
		return nil, err
	}
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = config.Endpoint
	producerConfig.AccessKeyID = config.AccessKey
	producerConfig.AccessKeySecret = config.SecretKey

	client = &Client{
		Config:   config,
		Producer: producer.InitProducer(producerConfig),
	}
	client.Producer.Start()
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{"content": "log-start"})
	err = client.Producer.SendLog(config.Project, config.LogStore, "start", config.HostName, logmsg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Info 记录info日志
func (c *Client) Info(topic string, logs any) error {
	logsMap := make(map[string]string)
	// type switch
	switch vt := logs.(type) {
	case map[string]string:
		logsMap = vt
	case map[string]any:
		for i, v := range vt {
			logsMap[i] = util.MarshalString(v)
		}
	case []byte:
		logsTmp := map[string]any{}
		err := json.Unmarshal(vt, &logsTmp)
		if err != nil {
			return ErrJsonStr
		}
		for i, v := range logsTmp {
			logsMap[i] = util.MarshalString(v)
		}
	default: // 结构体
		logsMap["content"] = util.MarshalString(logs)
	}
	return c.Record("info", topic, logsMap)
}

func (c *Client) Record(level string, topic string, logs map[string]string) error {
	ts := time.Now().Unix()
	logs["level"] = level
	logs["log_ts"] = util.Int642String(ts)
	// send
	slsLog := producer.GenerateLog(uint32(ts), logs)
	return c.Producer.SendLog(c.Config.Project, c.Config.LogStore, topic, c.Config.HostName, slsLog)
}

// Close 关闭日志服务
func (c *Client) Close() {
	if c.Producer != nil {
		c.Producer.SafeClose()
	}
}

// checkConfig 验证配置是否缺少 自动创建LogStore
func checkConfig(conf *Config) (err error) {
	if conf.AccessKey == "" || conf.Endpoint == "" || conf.Project == "" || conf.LogStore == "" || conf.SecretKey == "" {
		return ErrInvalidParam
	}
	return
}
