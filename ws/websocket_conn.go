package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

/*
TODO: 实现 gorilla/websocket 的事件监听，适配 socketio
socketio 数据传输格式： `["binding","{\"code\":0,\"data\":123456,\"message\":\"\"}"]`
// 字符串数组，第一个字符串是事件名，第二个字符串是数据，通常是 json string
["",""]
*/
type Conn struct {
	SocketId       string
	wsConn         *websocket.Conn
	inChan         chan []byte
	outJsonChan    chan any
	outMessageChan chan []byte
	closeChan      chan byte
	mutex          sync.Mutex
	isClosed       bool
}

func InitConn(wsConn *websocket.Conn) (conn *Conn, err error) {
	conn = &Conn{
		SocketId:       wsConn.RemoteAddr().String(),
		wsConn:         wsConn,
		inChan:         make(chan []byte, CommonChanCount),
		outJsonChan:    make(chan any, CommonChanCount),
		outMessageChan: make(chan []byte, CommonChanCount),
		closeChan:      make(chan byte, 1),
	}

	// 启动协程去读取消息
	go conn.readLoop()
	// 启动协程去发送消息
	go conn.writeLoop()
	return
}

func (c *Conn) ReadMessage() (data []byte, err error) {
	select {
	case data = <-c.inChan:
	case <-c.closeChan:
		err = ErrConnClosed
	}
	return data, err
}

func (c *Conn) WriteMessage(data []byte) (err error) {
	select {
	case c.outMessageChan <- data:
	case <-c.closeChan:
		err = ErrConnClosed
	}
	return err
}

func (c *Conn) WriteJson(data any) (err error) {
	select {
	case c.outJsonChan <- data:
	case <-c.closeChan:
		err = ErrConnClosed
	}
	return err
}

func (c *Conn) Close() {
	// 线程安全的Close
	_ = c.wsConn.Close()

	// 这一行代码只需要执行一次
	c.mutex.Lock()
	if !c.isClosed {
		close(c.closeChan)
		c.isClosed = true
	}
	c.mutex.Unlock()
}

func (c *Conn) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		_, data, err = c.wsConn.ReadMessage()
		if err != nil {
			goto ERROR
		}
		select {
		case c.inChan <- data:
		case <-c.closeChan:
			// todo:关闭
			goto ERROR
		}
	}

ERROR:
	// todo:关闭连接操作
	c.Close()
}

func (c *Conn) writeLoop() {
	var (
		message []byte
		json    any
		err     error
	)
	for {
		select {
		case message = <-c.outMessageChan:
			err = c.wsConn.WriteMessage(websocket.TextMessage, message)
		case json = <-c.outJsonChan:
			err = c.wsConn.WriteJSON(json)
		case <-c.closeChan:
			goto ERROR
		}

		if err != nil {
			goto ERROR
		}
	}
ERROR:
	//todo:关闭连接操作
	c.Close()
}
