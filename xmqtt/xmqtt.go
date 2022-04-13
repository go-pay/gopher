package xmqtt

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pay/gopher/retry"
	"github.com/go-pay/gopher/xlog"
	"github.com/google/uuid"
)

type Client struct {
	c        *Config
	mu       sync.Mutex
	Ops      *mqtt.ClientOptions
	Mqtt     mqtt.Client
	Topics   []string
	SubFuncs map[string]mqtt.MessageHandler // key:topic#qos, value: callback func
}

// New 1、New
func New(c *Config) (mc *Client) {
	// 日志
	mqtt.ERROR = log.New(os.Stderr, "[MQTT.ERROR] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	var (
		clientId = c.ClientId
		ops      = mqtt.NewClientOptions()
	)
	ops.AddBroker(fmt.Sprintf("tcp://%s:%d", c.Broker, c.TcpPort))
	if clientId == "" {
		clientId = uuid.NewString()
	}
	ops.SetClientID(clientId)
	ops.SetUsername(c.Uname)
	ops.SetPassword(c.Password)
	if c.KeepAlive > 0 {
		ops.SetKeepAlive(c.KeepAlive)
	}
	ops.SetCleanSession(c.CleanSession)
	mc = &Client{
		c:        c,
		Ops:      ops,
		SubFuncs: make(map[string]mqtt.MessageHandler),
	}
	return mc
}

// OnConnectListener 2、设置链接监听
func (c *Client) OnConnectListener(fun mqtt.OnConnectHandler) (mc *Client) {
	if fun != nil {
		c.Ops.OnConnect = fun
	}
	return c
}

// OnConnectLostListener 3、设置断开链接监听
func (c *Client) OnConnectLostListener(fun mqtt.ConnectionLostHandler) (mc *Client) {
	if fun != nil {
		c.Ops.OnConnectionLost = fun
	}
	return c
}

// StartAndConnect 4、真实创建Client并连接mqtt
func (c *Client) StartAndConnect() (err error) {
	if c.Ops.OnConnect == nil {
		c.Ops.OnConnect = c.DefaultOnConnectFunc
	}
	// new
	c.Mqtt = mqtt.NewClient(c.Ops)
	err = retry.Retry(func() error {
		token := c.Mqtt.Connect()
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
		return nil
	}, 3, 3*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Close 主动断开连接
func (c *Client) Close() {
	c.Mqtt.Disconnect(1000)
}

// Subscribe 订阅topic
func (c *Client) Subscribe(topic string, qos QosType, callback mqtt.MessageHandler) error {
	// callback 缓存，断开连接后，重新注册订阅
	c.mu.Lock()
	c.SubFuncs[subCallbackKey(topic, qos)] = callback
	c.mu.Unlock()
	if err := c.sub(topic, qos, callback); err != nil {
		return err
	}
	return nil
}

func (c *Client) sub(topic string, qos QosType, callback mqtt.MessageHandler) error {
	token := c.Mqtt.Subscribe(topic, byte(qos), callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// UnSubscribe 取消订阅topic
func (c *Client) UnSubscribe(topics ...string) error {
	token := c.Mqtt.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish 推送消息
func (c *Client) Publish(topic string, qos QosType, payload interface{}) error {
	token := c.Mqtt.Publish(topic, byte(qos), false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func subCallbackKey(topic string, qos QosType) string {
	return fmt.Sprintf("%s#%v", topic, qos)
}

func (c *Client) DefaultOnConnectFunc(cli mqtt.Client) {
	xlog.Infof("ClientId [%s] Connected", c.Ops.ClientID)
	c.mu.Lock()
	defer c.mu.Unlock()
	// 连接后，注册订阅
	for key, handler := range c.SubFuncs {
		// 协程
		go func(k string, cb mqtt.MessageHandler) {
			split := strings.Split(k, "#")
			if len(split) == 2 {
				var qos QosType
				switch split[1] {
				case "0":
					qos = QosAtMostOne
				case "1":
					qos = QosAtLeastOne
				case "2":
					qos = QosOnlyOne
				default:
					qos = QosAtMostOne
				}
				err := retry.Retry(func() error {
					return c.sub(split[0], qos, cb)
				}, 3, 2*time.Second)
				if err != nil {
					xlog.Errorf("topic[%s] sub callback register err:%+v", split[0], err)
				}
			}
		}(key, handler)
	}
}
