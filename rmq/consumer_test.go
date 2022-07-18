package rmq

import (
	"testing"
)

func TestConsumer(t *testing.T) {
	//conf := &RocketMQConfig{
	//	Namespace: "MQ_INST_xxx",
	//	GroupName: "GID_xxx",
	//	EndPoint:  "http://xxx.cn-hangzhou.mq-internal.aliyuncs.com:8080",
	//	AccessKey: "xxx",
	//	SecretKey: "xxx",
	//	LogLevel:  LogError,
	//	// 自定义配置
	//	//ConsumerOptions: []consumer.Option{
	//	//	consumer.WithMaxReconsumeTimes(10),
	//	//	consumer.WithConsumeMessageBatchMaxSize(1),
	//	//	consumer.WithPullInterval(time.Millisecond),
	//	//	consumer.WithPullBatchSize(10),
	//	//	consumer.WithConsumerModel(consumer.BroadCasting),
	//	//},
	//}
	//conn, err := NewConsumer(conf).Conn()
	//if err != nil {
	//	xlog.Error(err)
	//	return
	//}
	//defer conn.Close()
	//
	//if err = conn.SubscribeSingle("topic", "*", func(c context.Context, ext *primitive.MessageExt) error {
	//	xlog.Debugf("body:%v", string(ext.Body))
	//	return nil
	//}); err != nil {
	//	xlog.Error(err)
	//	return
	//}
	//
	//if err = conn.Start(); err != nil {
	//	xlog.Error(err)
	//}
	//
	//time.Sleep(time.Hour)
}
