package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pay/xlog"
	"github.com/go-pay/gopher/xmqtt"
)

func main2() {
	// 初始化参数和配置
	emqx := xmqtt.New(&xmqtt.Config{
		Broker:       "127.0.0.1",
		TcpPort:      1883,
		ClientId:     "clientid_producer",
		Uname:        "uname",
		Password:     "password",
		CleanSession: true,
	})
	// 设置Mqtt连接监听
	emqx.OnConnectListener(emqx.DefaultOnConnectFunc)
	// 设置Mqtt断开连接监听
	emqx.OnConnectLostListener(func(client mqtt.Client, err error) {
		xlog.Infof("[%s]IsConnected[%+t] lost connection, err: %+v", emqx.Ops.ClientID, client.IsConnected(), err)
	})
	// 启动连接
	if err := emqx.StartAndConnect(); err != nil {
		panic(err)
	}

	// 发布消息
	_ = emqx.Publish("topic1", xmqtt.QosAtMostOne, "hello world")

}
