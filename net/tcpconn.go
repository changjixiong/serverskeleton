package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type TCPConnection struct {
	SendBuf chan []byte
	ReadBuf chan []byte
	Conn    net.Conn
	client  *Client
}

func (c *TCPConnection) Run(addr string) {

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

			c.Conn = rw
			c.SendBuf = make(chan []byte, 10)
			c.ReadBuf = make(chan []byte, 10)
			c.client = &Client{}

			go c.read()
			go c.send()
		}
	}()
}
func (c *TCPConnection) read() {
	buf := make([]byte, 1024)
	defer c.Conn.Close()

	for {
		n, err := c.Conn.Read(buf)
		// fmt.Println("recv:", string(buf[0:n]))
		switch err {
		case nil:
			c.handle(buf[0:n])

		case io.EOF:
			fmt.Printf("Warning: End of data: %s \n", err)
			return
		default:
			fmt.Printf("Error: Reading data : %s \n", err)
			return
		}
	}
}

func (c *TCPConnection) handle(msg []byte) {
	c.SendBuf <- msg
}
func (c *TCPConnection) send() {
	for {
		respon := <-c.SendBuf
		data, _ := json.Marshal(&respon)
		c.Conn.Write(data)
	}

}
