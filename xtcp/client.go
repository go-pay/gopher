package xtcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/go-pay/gopher/util"
	"github.com/go-pay/gopher/xlog"
	"github.com/google/uuid"
	"io"
	"net"
	"strings"
	"time"
)

// TCPClient 客户端实例
type TCPClient struct {
	// 服务端连接的唯一 ID
	UID string
	// 连接句柄
	conn *FD
	// 是否注册了关闭方法
	haveRegisterClose bool
	// 是否注册了消息处理方法
	haveRegisterHandleMessage bool
	// 关闭事件管道
	closeChan chan *Context
	// 处理消息管道
	handleMessageChan chan *Context
}

// NewClient New TCP Client
func NewClient(c *Config) (*TCPClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	if c.ClientId == "" {
		c.ClientId = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	client := &TCPClient{
		UID:               c.ClientId, // wait server return
		conn:              &FD{conn: conn},
		closeChan:         make(chan *Context),
		handleMessageChan: make(chan *Context),
	}
	go func() {
		for {
			if err := client.conn.SendByte([]byte("PING")); err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					co, err := net.DialTCP("tcp", nil, addr)
					if err != nil {
						xlog.Errorf("TCPClient.reconnect: %+v", err)
						time.Sleep(c.HeartBeat)
						continue
					}
					client.conn.conn = co
				}
			}
			time.Sleep(c.HeartBeat)
		}
	}()
	return client, nil
}

// listen 客户端消息监听
func (c *TCPClient) listen() {
	reader := bufio.NewReader(c.conn.conn)
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorf("TCPClient.listen.Panic: %+v", err)
		}
	}()
	for {
		// 前4个字节表示数据长度
		// 此外 Peek 方法并不会减少 reader 中的实际数据量
		peek, err := reader.Peek(4)
		if err != nil {
			xlog.Color(xlog.Red).Errorf("TCPClient.listen.Peek: %+v", err)
			c.Close(err.Error())
			break
		}
		buffer := bytes.NewBuffer(peek)
		var length int32
		// 读取缓冲区前4位,代表消息实体的数据长度,赋予 length 变量
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			panic(err)
		}
		// reader.Buffered() 返回缓存中未读取的数据的长度,
		// 如果缓存区的数据小于总长度，则意味着数据不完整
		if int32(reader.Buffered()) < length+4 {
			continue
		}
		//从缓存区读取大小为数据长度的数据
		data := make([]byte, length+4)
		_, err = reader.Read(data)
		if err != nil {
			if err == io.EOF {
				c.Close(err.Error())
				break
			}
			continue
		}
		// 如果没有注册该方法,则丢弃这条消息,继续下一轮
		if !c.haveRegisterHandleMessage {
			continue
		}
		m := new(Message)
		//xlog.Infof("client.data[4:]，%v", string(data[4:]))
		_ = util.UnmarshalBytes(data[4:], m)
		// 管道分发，事件处理
		c.handleMessageChan <- &Context{
			uid:      c.UID,
			remoteIP: c.conn.conn.RemoteAddr().String(),
			conn:     c.conn,
			message:  m.Body,
		}
	}
}

// Run 启动方法
func (c *TCPClient) Run() {
	xlog.Warnf("[%s]Connecting TCP server on %s", c.UID, c.conn.conn.RemoteAddr().String())
	go c.listen()
}

// GracefulClose 优雅退出,该方法只允许被调用一回
func (c *TCPClient) Close(msg ...string) {
	// 关闭连接
	defer c.conn.close()             // 3
	defer close(c.closeChan)         // 2
	defer close(c.handleMessageChan) // 1
	// 事件通知
	if c.haveRegisterClose {
		message := ""
		if len(msg) > 0 {
			message = msg[0]
		}
		c.closeChan <- &Context{
			uid:      c.UID,
			remoteIP: c.conn.conn.RemoteAddr().String(),
			conn:     c.conn,
			message:  message,
		}
	}
}

// =====================================================================================================================

// OnConnect 建立连接触发的回调
func (c *TCPClient) OnConnect(f HandleFunc) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPClient.OnConnect.Panic: %+v", err)
			}
		}()
		// do func
		f(&Context{
			uid:      c.UID,
			remoteIP: c.conn.conn.RemoteAddr().String(),
			conn:     c.conn,
			message:  "",
		})
	}()
}

// OnMessage 当收到消息时触发
func (c *TCPClient) OnMessage(f HandleFunc) {
	c.haveRegisterHandleMessage = true
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPClient.OnConnect.Panic: %+v", err)
			}
		}()
		for ctx := range c.handleMessageChan {
			go f(ctx)
		}
	}()
}

// OnClose 连接断开触发的回调
func (c *TCPClient) OnClose(f HandleFunc) {
	c.haveRegisterClose = true
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPClient.OnClose.Panic: %+v", err)
			}
		}()
		for ctx := range c.closeChan {
			go func(ctx *Context) {
				defer func() {
					if err := recover(); err != nil {
						xlog.Errorf("TCPClient.OnClose.Handel.Panic: %+v", err)
					}
				}()
				// do func
				f(ctx)
			}(ctx)
		}
	}()
}

// =====================================================================================================================

// SendText 实现客户端发送消息给服务端的方法
func (c *TCPClient) SendText(msg string) error {
	return c.conn.SendText(msg)
}

// SendByte 发送一条消息
func (c *TCPClient) SendByte(msg []byte) error {
	return c.conn.SendByte(msg)
}
