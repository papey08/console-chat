package wsserver

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WsServer struct {
	connections map[string]net.Conn
	mu          *sync.Mutex
}

func NewWsServer() *WsServer {
	return &WsServer{
		connections: make(map[string]net.Conn),
		mu:          new(sync.Mutex),
	}
}

func (s *WsServer) AddConnection(conn net.Conn, nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[nickname] = conn
}

func (s *WsServer) Chat(w http.ResponseWriter, r *http.Request) {
	// creating connection to websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("can't upgrade connection:", err.Error())
		return
	}

	// getting info about user
	usrData, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		log.Println("can't read user data:", err.Error())
		return
	}
	var usr userPresenter
	if err := json.Unmarshal(usrData, &usr); err != nil {
		log.Println("can't unmarshal user data:", err.Error())
		return
	}

	// creating connection for new user
	s.AddConnection(conn, usr.Nickname)
	ch := make(chan []byte)

	// reading new messages
	go func() {
		defer func() {
			close(ch)
			if err := conn.Close(); err != nil {
				log.Fatal("connection closure error:", err.Error())
			}
		}()

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil && err != io.EOF {
				log.Fatal("can't read message from connection:", err.Error())
			} else if err == io.EOF {
				break
			}

			ch <- msg
		}
	}()

	// sending messages to the users
	go func() {
		for msg := range ch {
			msg = append([]byte(usr.Nickname+": "), msg...)
			s.mu.Lock()
			for key, connection := range s.connections {
				if key == usr.Nickname {
					continue
				} else if err := wsutil.WriteServerMessage(connection, ws.OpText, msg); err != nil {
					log.Println("can't write message:", err.Error())
					delete(s.connections, key)
				}
			}
			s.mu.Unlock()
		}
		log.Println(usr.Nickname, "leaves the chat")
		s.mu.Lock()
		delete(s.connections, usr.Nickname)
		s.mu.Unlock()
	}()
}
