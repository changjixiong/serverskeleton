package net

type Client struct {
}

type ServerManager struct {
}

func (s *ServerManager) ListenAndServe(ser Server, addr string) {
	ser.Run(addr)
}

type Server interface {
	Run(addr string)
}
