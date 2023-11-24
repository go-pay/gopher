package xmqtt

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pay/gopher/retry"
	"github.com/go-pay/xlog"
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
	if len(c.Topics) > 0 {
		_ = c.UnSubscribe(c.Topics...)
		c.Topics = nil
	}
	c.SubFuncs = nil
	c.Mqtt.Disconnect(1000)
}

func (c *Client) DefaultOnConnectFunc(cli mqtt.Client) {
	xlog.Infof("Clientid [%s] Connected", c.Ops.ClientID)
	c.mu.Lock()
	defer c.mu.Unlock()
	// 若 c.SubFuncs 不为空，连接后注册订阅
	for key, handler := range c.SubFuncs {
		// 协程
		go func(k string, cb mqtt.MessageHandler) {
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 64<<10)
					buf = buf[:runtime.Stack(buf, false)]
					xlog.Errorf("reSubscribe: panic recovered: %s\n%s", r, buf)
				}
			}()
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
