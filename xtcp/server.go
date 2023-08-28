package xtcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/go-pay/gopher/smap"
	"github.com/go-pay/gopher/util"
	"github.com/go-pay/gopher/xlog"
	"github.com/google/uuid"
	"net"
	"time"
)

type TCPServer struct {
	// 心跳
	heartBeat time.Duration
	// 服务端监听句柄
	listener net.Listener
	// 服务端文件描述符的上下文集合 key 是 uid 值为 Context
	fdsContext smap.Map[string, *Context]
	// 是否注册连接方法
	haveRegisterConn bool
	// 是否注册了关闭方法
	haveRegisterClose bool
	// 是否注册了消息处理方法
	haveRegisterHandleMessage bool
	// 连接事件管道
	connChan chan *Context
	// 关闭事件管道
	closeChan chan *Context
	// 处理消息管道
	handleMessageChan chan *Context
	// 是否关闭
	isShutdown bool
}

// NewServer New TCP Server
func NewServer(c *Config) (*TCPServer, error) {
	listener, err := net.Listen("tcp", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}
	if c.HeartBeat == 0 {
		c.HeartBeat = time.Second * 5
	}
	s := &TCPServer{
		listener:          listener,
		closeChan:         make(chan *Context),
		connChan:          make(chan *Context),
		handleMessageChan: make(chan *Context),
		heartBeat:         c.HeartBeat,
	}
	return s, nil
}

// accept 接受连接
func (s *TCPServer) accept() {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorf("TCPServer.accept.Panic: %+v", err)
		}
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil && s.isShutdown {
			break
		}
		xlog.Infof("accept new conn: %s", conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}

		uid := uuid.New().String()
		// 构建上下文
		ctx := &Context{
			uid:      uid,
			remoteIP: conn.RemoteAddr().String(),
			conn:     &FD{conn: conn},
		}
		// 将连接维护到总集合中
		s.fdsContext.Store(uid, ctx)

		// 如果用户注册了以下方法
		if s.haveRegisterConn {
			s.connChan <- ctx
			//if err := ctx.conn.SendText("CID_" + uid); err != nil {
			//	xlog.Errorf("TCPServer.SetClientCID.Err: %+v", err)
			//	continue
			//}
		}

		if s.haveRegisterHandleMessage {
			s.handleMessageChan <- ctx
		}
	}
}

// Run 启动方法
func (s *TCPServer) Run() {
	xlog.Warnf("Listening and serving TCP on %s", s.listener.Addr().String())
	go s.accept()
}

// CloseByUID 通过 uid 关闭一个连接
func (s *TCPServer) CloseByUID(uid string) {
	if fd, exists := s.fdsContext.Load(uid); exists {
		// 优先从集合中删除,再关闭连接,请注意执行顺序不可颠倒
		s.fdsContext.Delete(uid)
		// 关闭连接
		fd.conn.close()
		// 推送关闭事件
		if s.haveRegisterClose {
			s.closeChan <- fd
		}
	}
}

// =====================================================================================================================

// BroadcastText 广播
func (s *TCPServer) BroadcastText(msg string) {
	s.fdsContext.Range(func(key string, ctx *Context) bool {
		if err := ctx.SendText(msg); err != nil {
			xlog.Errorf("BroadcastText.Err: %+v", err)
		}
		return true
	})
}

// BroadcastTextOther 广播给其它人,除了自己
func (s *TCPServer) BroadcastTextOther(uid string, msg string) {
	s.fdsContext.Range(func(key string, ctx *Context) bool {
		if id := ctx.uid; uid == id {
			return true
		}
		if err := ctx.SendText(msg); err != nil {
			xlog.Errorf("BroadcastTextOther.Err: %+v", err)
		}
		return true
	})
}

// BroadcastByte 广播
func (s *TCPServer) BroadcastByte(msg []byte) {
	s.fdsContext.Range(func(key string, ctx *Context) bool {
		if err := ctx.SendByte(msg); err != nil {
			xlog.Errorf("BroadcastByte.Err: %+v", err)
		}
		return true
	})
}

// BroadcastByteOther 广播给其它人,除了自己
func (s *TCPServer) BroadcastByteOther(uid string, msg []byte) {
	s.fdsContext.Range(func(key string, ctx *Context) bool {
		if id := ctx.uid; uid == id {
			return true
		}
		if err := ctx.SendByte(msg); err != nil {
			xlog.Errorf("BroadcastByteOther.Err: %+v", err)
		}
		return true
	})
}

// SendText 发送一条消息
func (s *TCPServer) SendText(uid string, msg string) error {
	if ctx, ok := s.fdsContext.Load(uid); ok {
		return ctx.SendText(msg)
	}
	return nil
}

// SendByte 发送一条消息
func (s *TCPServer) SendByte(uid string, msg []byte) error {
	if ctx, ok := s.fdsContext.Load(uid); ok {
		return ctx.SendByte(msg)
	}
	return nil
}

// Close 退出
func (s *TCPServer) Close() {
	if s.isShutdown {
		return
	}
	s.isShutdown = true
	_ = s.listener.Close()
	defer close(s.closeChan)         // 3
	defer close(s.handleMessageChan) // 2
	defer close(s.connChan)          // 1
	s.fdsContext.Range(func(key string, v *Context) bool {
		s.CloseByUID(key)
		return true
	})
}

// =====================================================================================================================

// OnConnect 当有连接发生时
func (s *TCPServer) OnConnect(f HandleFunc) {
	s.haveRegisterConn = true
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPServer.OnConnect.Panic: %+v", err)
			}
		}()
		for ctx := range s.connChan {
			go func(ctx *Context) {
				defer func() {
					if err := recover(); err != nil {
						xlog.Errorf("TCPServer.OnConnect.Handel.Panic: %+v", err)
					}
				}()
				// do func
				f(ctx)
			}(ctx)
		}
	}()
}

// OnMessage 当有消息时
func (s *TCPServer) OnMessage(f HandleFunc) {
	s.haveRegisterHandleMessage = true
	// 注册事件
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPServer.OnMessage.Panic: %+v", err)
			}
		}()
		for ctx := range s.handleMessageChan {
			// 消息并发处理
			go func(ctx *Context) {
				defer func() {
					if err := recover(); err != nil {
						xlog.Errorf("TCPServer.OnMessage.Handel.Panic: %+v", err)
					}
				}()
				fd := ctx.Conn()
				reader := bufio.NewReader(fd.conn)
				for {
					err := fd.conn.SetDeadline(time.Now().Add(s.heartBeat))
					if err != nil {
						xlog.Warnf("conn.SetDeadline.Err: %+v", err)
						err = nil
					}
					// 前4个字节表示数据长度
					// 此外 Peek 方法并不会减少 reader 中的实际数据量
					peeked, err := reader.Peek(4)
					if err != nil {
						// 读取失败关闭连接
						s.CloseByUID(ctx.uid)
						break
					}
					buffer := bytes.NewBuffer(peeked)
					var length int32
					// 读取缓冲区前4位,代表消息实体的数据长度,赋予 length 变量
					err = binary.Read(buffer, binary.BigEndian, &length)
					if err != nil {
						continue
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
						continue
					}
					m := new(Message)
					//xlog.Infof("server.data[4:]，%v", string(data[4:]))
					_ = util.UnmarshalBytes(data[4:], m)
					// 将消息内容赋给上下文
					ctx.message = m.Body
					// do func
					f(ctx)
				}
			}(ctx)
		}
	}()
}

// OnClose 当客户端断开时
func (s *TCPServer) OnClose(f HandleFunc) {
	s.haveRegisterClose = true
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf("TCPServer.OnClose.Panic: %+v", err)
			}
		}()
		for ctx := range s.closeChan {
			go func(ctx *Context) {
				defer func() {
					if err := recover(); err != nil {
						xlog.Errorf("TCPServer.OnClose.Handel.Panic: %+v", err)
					}
				}()
				// do func
				f(ctx)
			}(ctx)
		}
	}()
}
