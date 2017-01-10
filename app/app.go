package main

import (
	"fmt"
	"os"
	"os/signal"
	"serverskeleton/net"
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

	tcpServer := &net.TPCServer{}
	wsServer := &net.WSServer{}

	server := &net.ServerManager{}
	server.ListenAndServe(tcpServer, ":8090")
	server.ListenAndServe(wsServer, ":8080")

	handleSignal()
}
