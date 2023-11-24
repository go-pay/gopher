package rmq

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-pay/limiter"
)

const (
	LogDebug LogLevel = "debug"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
	LogInfo  LogLevel = "info"
)

type LogLevel string

type RocketMQConfig struct {
	// 阿里云 实例ID
	Namespace string
	// GroupID 阿里云创建
	GroupName string
	// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。
	EndPoint string
	// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
	AccessKey string
	// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
	SecretKey string
	// log 级别 // default info
	LogLevel LogLevel
	// 自定义消费者配置
	ConsumerOptions []consumer.Option
	// 自定义生产者配置
	ProducerOptions []producer.Option
	// currently consume limiter
	Limit *limiter.Config
}

func defaultConsumerOps(conf *RocketMQConfig) (ops []consumer.Option) {
	ops = []consumer.Option{
		consumer.WithNamespace(conf.Namespace),
		consumer.WithGroupName(conf.GroupName),
		consumer.WithNameServer(primitive.NamesrvAddr{conf.EndPoint}),
		consumer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		consumer.WithRetry(3),
	}
	return ops
}

func defaultProducerOps(conf *RocketMQConfig) (ops []producer.Option) {
	ops = []producer.Option{
		producer.WithNamespace(conf.Namespace),
		producer.WithNameServer(primitive.NamesrvAddr{conf.EndPoint}),
		producer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		producer.WithRetry(3),
	}
	// GroupName is not necessary for producer
	if conf.GroupName != "" {
		ops = append(ops, producer.WithGroupName(conf.GroupName))
	}
	return ops
}
