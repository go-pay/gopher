package ws

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	SocketId       string
	wsConn         *websocket.Conn
	inChan         chan []byte
	outJsonChan    chan interface{}
	outMessageChan chan []byte
	closeChan      chan byte
	mutex          sync.Mutex
	isClosed       bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		SocketId:       wsConn.RemoteAddr().String(),
		wsConn:         wsConn,
		inChan:         make(chan []byte, 1000),
		outJsonChan:    make(chan interface{}, 1000),
		outMessageChan: make(chan []byte, 1000),
		closeChan:      make(chan byte, 1),
	}

	//启动协程去读取消息
	go conn.readLoop()
	//启动协程去发送消息
	go conn.writeLoop()
	return
}

func (this *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-this.inChan:
	case <-this.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (this *Connection) WriteMessage(data []byte) (err error) {
	select {
	case this.outMessageChan <- data:
	case <-this.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (this *Connection) WriteJson(data interface{}) (err error) {
	select {
	case this.outJsonChan <- data:
	case <-this.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (this *Connection) Close() {
	//线程安全的Close
	this.wsConn.Close()

	//这一行代码只需要执行一次
	this.mutex.Lock()
	if !this.isClosed {
		close(this.closeChan)
		this.isClosed = true
	}
	this.mutex.Unlock()
}

func (this *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		_, data, err = this.wsConn.ReadMessage()
		if err != nil {
			goto ERROR
		}
		select {
		case this.inChan <- data:
		case <-this.closeChan:
			//todo:关闭
			goto ERROR
		}
	}

ERROR:
	//todo:关闭连接操作
	this.Close()
}

func (this *Connection) writeLoop() {
	var (
		message []byte
		json    interface{}
		err     error
	)
	for {
		select {
		case message = <-this.outMessageChan:
			err = this.wsConn.WriteMessage(websocket.TextMessage, message)
		case json = <-this.outJsonChan:
			err = this.wsConn.WriteJSON(json)
		case <-this.closeChan:
			goto ERROR
		}

		if err != nil {
			goto ERROR
		}
	}
ERROR:
	//todo:关闭连接操作
	this.Close()
}
