package wsserver

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

// codeNicknameInToken codes user nickname into valid token []byte
func codeNicknameInToken(nickname string) ([]byte, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["nickname"] = nickname
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	tokenInStr, err := token.SignedString([]byte("abcd"))
	if err != nil {
		return nil, err
	}
	return []byte(tokenInStr), err
}

// getChat connects to the chat and sends user token to authorize
func getChat(url string, token []byte) (net.Conn, error) {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		return nil, err
	}
	if err = wsutil.WriteClientMessage(conn, ws.OpText, token); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func TestChat(t *testing.T) {
	wsserver := NewWsServer([]byte("abcd"))
	server := httptest.NewServer(http.HandlerFunc(wsserver.Chat))
	defer server.Close()

	url := "ws" + server.URL[4:]

	// user01 joins the chat
	token, err := codeNicknameInToken("user01")
	assert.NoError(t, err)
	conn01, err := getChat(url, token)
	assert.NoError(t, err)
	defer conn01.Close()

	// user02 joins the chat
	token, err = codeNicknameInToken("user02")
	assert.NoError(t, err)
	conn02, err := getChat(url, token)
	assert.NoError(t, err)
	defer conn02.Close()

	// these and all following sleeps are to prevent
	// message race in the chat for correct testing
	time.Sleep(100 * time.Millisecond)

	// user01 gets message that user02 joined the chat
	msg, _, err := wsutil.ReadServerData(conn01)
	assert.Equal(t, "user02 joins the chat", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user01 sends message to the chat
	err = wsutil.WriteClientMessage(conn01, ws.OpText, []byte("Ping"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user02 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(t, "user01: Ping", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user02 sends message to user01
	err = wsutil.WriteClientMessage(conn02, ws.OpText, []byte("Pong"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user01 gets message from user02
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(t, "user02: Pong", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user03 joins the chat
	token, err = codeNicknameInToken("user03")
	assert.NoError(t, err)
	conn03, err := getChat(url, token)
	assert.NoError(t, err)
	defer conn03.Close()

	time.Sleep(100 * time.Millisecond)

	// user01 gets message that user03 joined the chat
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(t, "user03 joins the chat", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user02 gets message that user03 joined the chat
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(t, "user03 joins the chat", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user01 sends message to the chat
	err = wsutil.WriteClientMessage(conn01, ws.OpText, []byte("Ping 2"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user02 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(t, "user01: Ping 2", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user03 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn03)
	assert.Equal(t, "user01: Ping 2", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user02 sends message to user01
	err = wsutil.WriteClientMessage(conn02, ws.OpText, []byte("Pong 2"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user01 gets message from user02
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(t, "user02: Pong 2", string(msg))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user03 sends message to user01
	err = wsutil.WriteClientMessage(conn03, ws.OpText, []byte("Pong 2"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// user01 gets message from user03
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(t, "user03: Pong 2", string(msg))
	assert.NoError(t, err)
}
