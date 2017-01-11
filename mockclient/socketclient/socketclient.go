package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("net.Dial 127.0.0.1:8090 :", err)
		os.Exit(1)
	}

	jsonStr := `
								{
								    "func_name":"Login",
								    "params":[
								        "userName",
								        "password"
								    ]
								}
								`

	for {
		fmt.Fprintf(conn, jsonStr)
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		switch err {
		case nil:
			fmt.Printf("recv %d:%s\n", n, string(buf))
		case io.EOF:
			fmt.Println("detected closed LAN connection")
			return
		default:
			fmt.Println("conn.Read", err)
			return
		}

		time.Sleep(time.Second * 2)
	}

}
