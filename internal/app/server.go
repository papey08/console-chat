package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WsServer struct {
	connections map[uint64]net.Conn
	mu          *sync.Mutex
	index       uint64
}

func (s *WsServer) addConnection(conn net.Conn) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index
	s.connections[s.index] = conn
	s.index += 1
	return i
}

func (s *WsServer) chat(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("can't upgrade connection: %s\n", err.Error())
		return
	}

	connID := s.addConnection(conn)
	id := fmt.Sprintf("%d", connID)
	ch := make(chan []byte)

	go func() {
		defer func() {
			conn.Close()
			close(ch)
		}()

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Printf("can't read message from connection: %s\n", err.Error())
				}

				break
			}

			ch <- msg
		}
	}()

	go func() {
		for msg := range ch {
			msg = append([]byte(id+": "), msg...)
			s.mu.Lock()
			for key, connection := range s.connections {
				if err := wsutil.WriteServerMessage(connection, ws.OpText, msg); err != nil {
					log.Printf("can't write message: %s\n", err.Error())
					delete(s.connections, key)
				}
			}

			s.mu.Unlock()
		}
		log.Println("go func stop")
		s.mu.Lock()
		delete(s.connections, connID)
		s.mu.Unlock()
	}()
}
