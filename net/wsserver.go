package net

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSServer struct {
}

type WSConnection struct {
	SendBuf     chan []byte
	ReadBuf     chan []byte
	Conn        *websocket.Conn
	messageType int
	client      *Client
}

var upgrader = websocket.Upgrader{} // use default options
func (ws *WSServer) serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	conn := &WSConnection{
		Conn:    c,
		SendBuf: make(chan []byte, 10),
		ReadBuf: make(chan []byte, 10),
		client:  &Client{},
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

func (ws *WSConnection) read() {
	defer ws.Conn.Close()

	for {

		mt, message, err := ws.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		ws.messageType = mt

		ws.handle(message)

	}
}

func (ws *WSConnection) handle(msg []byte) {
	ws.SendBuf <- msg
}

func (ws *WSConnection) send() {
	for {

		ws.Conn.WriteMessage(ws.messageType, <-ws.SendBuf)
	}

}
