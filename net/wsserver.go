package net

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serverskeleton/parser"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	MethodMap map[string]*parser.MethodInfo
}

func (w *WSServer) RegisterMethod(v interface{}) {

	parser.RegisterMethod(w.MethodMap, v)

}

type WSConnection struct {
	SendBuf     chan *parser.Response
	ReadBuf     chan []byte
	Conn        *websocket.Conn
	messageType int
	client      *Client
	isClosed    bool
	MethodMap   map[string]*parser.MethodInfo
}

var upgrader = websocket.Upgrader{} // use default options
func (ws *WSServer) serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	conn := &WSConnection{
		Conn:      c,
		SendBuf:   make(chan *parser.Response, 10),
		ReadBuf:   make(chan []byte, 10),
		client:    &Client{},
		MethodMap: ws.MethodMap,
	}

	go conn.read()
	go conn.send()

}

func (ws *WSServer) Run(Addr string) {

	mux := &http.ServeMux{}

	mux.HandleFunc("/serve", ws.serve)
	server := &http.Server{Addr: Addr, Handler: mux}
	server.ListenAndServe()

}

func (ws *WSConnection) translate(data []byte) {
	req, ok := parser.GenRequest(data)

	if !ok {
		ws.write(parser.GenErrRespones(parser.ErrorCode_JsonError))

	} else {
		ws.handle(req)
	}

}

func (ws *WSConnection) read() {
	defer ws.close()

	for {

		if ws.isClosed {
			return
		}

		mt, message, err := ws.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		ws.messageType = mt

		ws.translate(message)

	}
}

func (ws *WSConnection) handle(req *parser.Request) {

	resp := parser.Invoke(ws.MethodMap, req, ws.client)

	ws.SendBuf <- resp

}

func (ws *WSConnection) write(resp *parser.Response) {
	ws.SendBuf <- resp
}

func (ws *WSConnection) close() {

	close(ws.ReadBuf)
	close(ws.SendBuf)
	ws.Conn.Close()

}

func (ws *WSConnection) send() {
	defer func() {
		ws.isClosed = true
	}()

	for {

		respon, ok := <-ws.SendBuf
		if !ok {
			ws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		data, _ := json.Marshal(&respon)

		err := ws.Conn.WriteMessage(ws.messageType, data)

		if nil != err {
			fmt.Println("ws.Conn.WriteMessage:", err)
			return
		}

		if respon.FuncName == "Error" {
			return
		}
	}

}
