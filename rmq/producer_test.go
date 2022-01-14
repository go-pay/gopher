package rmq

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestProducer(t *testing.T) {
	//conf := &RocketMQConfig{
	//	Namespace: "MQ_INST_xxx",
	//	GroupName: "GID_xxx",
	//	EndPoint:  "http://xxx.cn-hangzhou.mq-internal.aliyuncs.com:8080",
	//	AccessKey: "xxx",
	//	SecretKey: "xxx",
	//	LogLevel:  LogError,
	//}
	//
	//conn, err := NewProducer(conf).Conn()
	//if err != nil {
	//	xlog.Error(err)
	//	return
	//}
	//
	//defer conn.Close()
	//
	//for i := 0; i < 1; i++ {
	//	message := &primitive.Message{
	//		Topic: "topic",
	//		Body:  []byte("12346"),
	//	}
	//	message.WithTag("tag").
	//		WithKeys([]string{"key"}).
	//		WithProperty("__STARTDELIVERTIME", "1642150538000") // 阿里云 定时消息 属性，value 为 毫秒时间戳
	//
	//	/*res,*/
	//	err = conn.SendAsyncSingle(ctx, nil, message)
	//	if err != nil {
	//		xlog.Errorf("%v", err)
	//		return
	//	}
	//	//xlog.Debugf("res:%#v", res)
	//}
	//
	//time.Sleep(time.Hour)
}
