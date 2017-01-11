package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"serverskeleton/parser"
	"time"
)

type TPCServer struct {
	MethodMap map[string]*parser.MethodInfo
}

func (t *TPCServer) RegisterMethod(v interface{}) {

	parser.RegisterMethod(t.MethodMap, v)

}

type TCPConnection struct {
	SendBuf      chan *parser.Response
	ReadBuf      chan []byte
	Conn         net.Conn
	messageType  int
	client       *Client
	isClosed     bool
	MethodMap    map[string]*parser.MethodInfo
	sendTryTimes int
}

func (t *TPCServer) Run(addr string) {

	go func() {
		ln, err := net.Listen("tcp", addr)

		if nil != err {
			fmt.Println("err:", err)
			return
		}

		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			rw, err := ln.Accept()
			if err != nil {
				//这个err 是一个 OpError类型，实现了net.Error接口
				//Temporary 临时错误，可重试

				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					fmt.Printf("Accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}

				fmt.Println("err:", err)
				return
			}

			tempDelay = 0

			conn := &TCPConnection{
				Conn:         rw,
				SendBuf:      make(chan *parser.Response, 10),
				ReadBuf:      make(chan []byte, 10),
				client:       &Client{},
				MethodMap:    t.MethodMap,
				sendTryTimes: 10,
			}

			go conn.read()
			go conn.send()

		}
	}()
}

func (t *TCPConnection) translate(data []byte) {
	req, ok := parser.GenRequest(data)

	if !ok {
		t.write(parser.GenErrRespones(parser.ErrorCode_JsonError))

	} else {
		t.handle(req)
	}

}

func (t *TCPConnection) close() {

	close(t.ReadBuf)
	close(t.SendBuf)
	t.Conn.Close()

}

func (t *TCPConnection) read() {
	defer t.close()
	buf := make([]byte, 1024)

	for {
		if t.isClosed {
			return
		}

		n, err := t.Conn.Read(buf)
		switch err {
		case nil:
			// c.handle(buf[0:n])
			t.translate(buf[0:n])

		case io.EOF:
			fmt.Printf("Warning: End of data: %s \n", err)
			return
		default:
			fmt.Printf("Error: Reading data : %s \n", err)
			return
		}
	}
}

func (t *TCPConnection) handle(req *parser.Request) {

	resp := parser.Invoke(t.MethodMap, req, t.client)

	t.SendBuf <- resp

}

func (t *TCPConnection) write(resp *parser.Response) {
	t.SendBuf <- resp
}
func (t *TCPConnection) send() {

	defer func() {
		t.isClosed = true
	}()

	for {

		respon, ok := <-t.SendBuf
		if !ok {
			// c.Conn.Write(kickclient)
			return
		}

		data, _ := json.Marshal(&respon)

		haveTryTimes := 0
		for len(data) > 0 && haveTryTimes < t.sendTryTimes {
			n, err := t.Conn.Write(data)

			if nil != err {
				fmt.Println("t.Conn.Write:", err)
				return
			}
			data = data[n:]
			haveTryTimes += 1
		}

		if haveTryTimes >= t.sendTryTimes {
			return
		}
		if respon.FuncName == "Error" {
			return
		}

	}

}
