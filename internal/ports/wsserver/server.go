package wsserver

import (
	"errors"
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

// AddConnection adds new client connection to the server
func (s *WsServer) AddConnection(conn net.Conn, nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[nickname] = conn
}

// getToken reads client token
func (s *WsServer) getToken(conn net.Conn) ([]byte, error) {
	tokenData, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

// auth checks if token is valid and returns nickname coded in token
func (s *WsServer) auth(tokenData []byte) (string, error) {
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
		return "", errors.New("auth failure")
	}
}

// sendMessageToUsers sends message from publisher to all other clients
func (s *WsServer) sendMessageToUsers(publisher string, message []byte) {
	s.mu.Lock()
	for key, connection := range s.connections {
		if key == publisher {
			continue
		} else if err := wsutil.WriteServerMessage(connection, ws.OpText, message); err != nil {
			log.Println(publisher, "was disconnected from the chat")
			delete(s.connections, key)
		}
	}
	defer s.mu.Unlock()

}

// Chat adds new client to the chat
func (s *WsServer) Chat(w http.ResponseWriter, r *http.Request) {
	// creating connection to websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("can't upgrade connection:", err.Error())
		return
	}

	// getting client's token and check if it is valid
	tokenData, err := s.getToken(conn)
	if err != nil {
		log.Println("can't upgrade connection:", err.Error())
		return
	}
	nickname, err := s.auth(tokenData)
	if err != nil {
		log.Println("can't get and validate token:", err.Error())
		return
	}

	// creating connection for new user
	s.AddConnection(conn, nickname)
	log.Println(nickname, "joins the chat")
	s.sendMessageToUsers(nickname, []byte(nickname+" joins the chat"))
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
			if err != nil {
				break
			}

			ch <- msg
		}
	}()

	// sending messages to the users
	go func() {
		for msg := range ch {
			msg = append([]byte(nickname+": "), msg...)
			s.sendMessageToUsers(nickname, msg)
		}
		log.Println(nickname, "leaves the chat")
		s.sendMessageToUsers(nickname, []byte(nickname+" leaves the chat"))
		s.mu.Lock()
		delete(s.connections, nickname)
		s.mu.Unlock()
	}()
}
