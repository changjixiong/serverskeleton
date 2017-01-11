package module

import (
	"fmt"

	"serverskeleton/net"
)

type PlayerModule struct {
	onlineClient map[int]*net.Client
}

var DefaultPlayerModule = &PlayerModule{}

func (p *PlayerModule) Login(client *net.Client, userName, pwd string) (data map[string]interface{}) {
	fmt.Println(client.Name, userName, pwd)

	return data

}
