package xtcp

import (
	"bytes"
	"encoding/binary"
	"github.com/go-pay/gopher/util"
	"net"
	"time"
)

// Config 配置
type Config struct {
	Host      string
	Port      string
	ClientId  string
	HeartBeat time.Duration
}

type HandleFunc func(ctx *Context)

// Context 上下文接收器
type Context struct {
	// 该连接的唯一 uid
	uid string
	// IP
	remoteIP string
	// 句柄
	conn *FD
	// 消息内容
	message string
}

// String 获取消息字符串
func (c *Context) Message() string {
	return c.message
}

// Byte 获取消息
func (c *Context) Byte() []byte {
	return []byte(c.message)
}

// RemoteIP 获取远程客户端信息
func (c *Context) RemoteIP() string {
	return c.remoteIP
}

// ConnUID 获取该连接唯一的 uid
func (c *Context) ConnUID() string {
	return c.uid
}

// Conn 获取连接句柄
func (c *Context) Conn() *FD {
	return c.conn
}

// SendText 发送消息
func (c *Context) SendText(msg string) error {
	return c.conn.SendText(msg)
}

// SendByte 发送字节
func (c *Context) SendByte(msg []byte) error {
	return c.conn.SendByte(msg)
}

// =====================================================================================================================

// Message 是一条标准消息的实现
type Message struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

// encode 消息编码
func (m *Message) encode() ([]byte, error) {
	// 序列化为 json
	msgBs := util.MarshalBytes(m)
	// 读取该 json 的长度
	var length = int32(len(msgBs))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.BigEndian, msgBs)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// FD 连接描述符的具体实现
type FD struct {
	conn net.Conn
}

// SendText 发送消息
func (c *FD) SendText(msg string) error {
	m := &Message{Body: msg}
	body, err := m.encode()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(body)
	return err
}

// SendByte 发送字节
func (c *FD) SendByte(msg []byte) error {
	m := &Message{Body: string(msg)}
	body, err := m.encode()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(body)
	return err
}

// GracefulClose 关闭文件描述符
func (c *FD) close() {
	_ = c.conn.Close()
}
