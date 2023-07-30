package wsserver

import (
	"net"
	"net/http"
	"sync"
)

type WsServer interface {
	Chat(w http.ResponseWriter, r *http.Request)
}

func New(tokenKey []byte) WsServer {
	return &wsServer{
		connections: make(map[string]net.Conn),
		mu:          new(sync.Mutex),
		tokenKey:    tokenKey,
	}
}
