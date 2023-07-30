package wsserver

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang-jwt/jwt"
)

type WsServer struct {
	connections map[string]net.Conn
	mu          *sync.Mutex
	tokenKey    []byte
}

func NewWsServer(tokenKey []byte) *WsServer {
	return &WsServer{
		connections: make(map[string]net.Conn),
		mu:          new(sync.Mutex),
		tokenKey:    tokenKey,
	}
}

func (s *WsServer) AddConnection(conn net.Conn, nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[nickname] = conn
}

func (s *WsServer) auth(conn net.Conn) (string, error) {
	tokenData, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(string(tokenData), func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.tokenKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["nickname"].(string), nil
	} else {
		return "", err
	}
}

func (s *WsServer) Chat(w http.ResponseWriter, r *http.Request) {
	// creating connection to websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("can't upgrade connection:", err.Error())
		return
	}

	// getting client's token and check if it is valid
	nickname, err := s.auth(conn)
	if err != nil {
		log.Println("can't get and validate token:", err.Error())
		return
	}

	// creating connection for new user
	s.AddConnection(conn, nickname)
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
			msg = append([]byte(nickname+": "), msg...)
			s.mu.Lock()
			for key, connection := range s.connections {
				if key == nickname {
					continue
				} else if err := wsutil.WriteServerMessage(connection, ws.OpText, msg); err != nil {
					log.Println("can't write message:", err.Error())
					delete(s.connections, key)
				}
			}
			s.mu.Unlock()
		}
		log.Println(nickname, "leaves the chat")
		s.mu.Lock()
		delete(s.connections, nickname)
		s.mu.Unlock()
	}()
}
