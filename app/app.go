package main

import (
	"fmt"
	"os"
	"os/signal"
	"serverskeleton/module"
	"serverskeleton/net"
	"serverskeleton/parser"
	"syscall"
)

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT,
		//syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTERM)
	for s := range c {
		fmt.Printf("Signals  get  %s ........... \n", s)
		break
	}
	fmt.Println("handleSignal end")

}

func main() {

	tcpServer := &net.TPCServer{MethodMap: make(map[string]*parser.MethodInfo)}
	tcpServer.RegisterMethod(module.DefaultPlayerModule)

	wsServer := &net.WSServer{MethodMap: make(map[string]*parser.MethodInfo)}
	wsServer.RegisterMethod(module.DefaultPlayerModule)

	server := &net.ServerManager{}
	server.ListenAndServe(tcpServer, ":8090")
	server.ListenAndServe(wsServer, ":8080")

	handleSignal()
}
