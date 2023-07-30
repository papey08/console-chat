package ginserver

import (
	"bytes"
	"console-chat/internal/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type tokenData struct {
	Data tokenResponse `json:"data"`
}

func (s *ginServerTestSuite) getUser(url string, body map[string]any) (tokenData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return tokenData{}, 0, err
	}
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/console-chat/users"+url, bytes.NewReader(data))
	if err != nil {
		return tokenData{}, 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	var resp tokenData
	code, err := s.getResponse(req, &resp)
	if err != nil {
		return tokenData{}, 0, err
	}
	return resp, code, nil
}

type getUserMock struct {
	nickname string
	password string
	usr      model.User
	err      error
}

func getHash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}

type getUserTest struct {
	description        string
	givenURL           string
	givenBody          map[string]any
	expectedStatusCode int
	expectGetToken     bool
	expectedNickname   string // to check if token is correct
}

func (s *ginServerTestSuite) TestGetUser() {
	mocks := []getUserMock{
		{
			nickname: "papey08",
			password: "qwerty_123",
			usr: model.User{
				Nickname:       "papey08",
				HashedPassword: "",
			},
			err: nil,
		},
		{
			nickname: "abcdefg",
			password: "qwerty_123",
			usr: model.User{
				Nickname:       "",
				HashedPassword: "",
			},
			err: model.UserNotFound,
		},
		{
			nickname: "papey09",
			password: "qwerty_123",
			usr: model.User{
				Nickname:       "",
				HashedPassword: "",
			},
			err: model.UserWrongPassword,
		},
	}
	tests := []getUserTest{
		{
			description: "correct signing in",
			givenURL:    "/papey08",
			givenBody: map[string]any{
				"password": "qwerty_123",
			},
			expectedStatusCode: http.StatusOK,
			expectGetToken:     true,
			expectedNickname:   "papey08",
		},
		{
			description: "user not exists",
			givenURL:    "/abcdefg",
			givenBody: map[string]any{
				"password": "qwerty_123",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description: "wrong password",
			givenURL:    "/papey09",
			givenBody: map[string]any{
				"password": "qwerty_123",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, m := range mocks {
		if m.usr.Nickname != "" {
			m.usr.HashedPassword = getHash(m.password)
		}
		s.app.On("SignInUser", mock.Anything, m.nickname, m.password).Return(m.usr, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.getUser(test.givenURL, test.givenBody)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)

		if test.expectGetToken {
			tokenStr := resp.Data.TokenString
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte("abcd"), nil
			})
			assert.NoError(s.T(), err)

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				assert.Equal(s.T(), test.expectedNickname, claims["nickname"].(string))
			} else {
				assert.Fail(s.T(), "wrong token")
			}
		}
	}
}

type userData struct {
	UserResp userResponse `json:"data"`
}

func (s *ginServerTestSuite) postUser(body map[string]any) (userData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return userData{}, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/console-chat/users", bytes.NewReader(data))
	if err != nil {
		return userData{}, 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	var resp userData
	code, err := s.getResponse(req, &resp)
	if err != nil {
		return userData{}, 0, err
	}
	return resp, code, nil
}

type postUserMock struct {
	nickname string
	password string
	usr      model.User
	err      error
}

type postUserTest struct {
	description        string
	givenBody          map[string]any
	expectedStatusCode int
}

func (s *ginServerTestSuite) TestPostUser() {
	mocks := []postUserMock{
		{
			nickname: "papey08",
			password: "qwerty_123",
			usr: model.User{
				Nickname:       "papey08",
				HashedPassword: "",
			},
			err: nil,
		},
		{
			nickname: "zzz",
			password: "qwerty_123",
			usr:      model.User{},
			err:      model.UserInvalidNickname,
		},
		{
			nickname: "papey08",
			password: "qwerty123",
			usr:      model.User{},
			err:      model.UserInvalidPassword,
		},
	}
	tests := []postUserTest{
		{
			description: "successful registration",
			givenBody: map[string]any{
				"nickname": "papey08",
				"password": "qwerty_123",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "invalid nickname",
			givenBody: map[string]any{
				"nickname": "zzz",
				"password": "qwerty_123",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "invalid password",
			givenBody: map[string]any{
				"nickname": "papey08",
				"password": "qwerty123",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("RegisterUser", mock.Anything, m.nickname, m.password).Return(m.usr, m.err).Once()
	}

	for _, test := range tests {
		_, code, err := s.postUser(test.givenBody)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}
}

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
func getChat(token []byte) (net.Conn, error) {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://localhost:8081/console-chat/chat")
	if err != nil {
		return nil, err
	}
	if err = wsutil.WriteClientMessage(conn, ws.OpText, token); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func (s *ginServerTestSuite) TestChat() {
	// user01 joins the chat
	token, err := codeNicknameInToken("user01")
	assert.NoError(s.T(), err)
	conn01, err := getChat(token)
	assert.NoError(s.T(), err)

	// user02 joins the chat
	token, err = codeNicknameInToken("user02")
	assert.NoError(s.T(), err)
	conn02, err := getChat(token)
	assert.NoError(s.T(), err)

	// user01 gets message that user02 joined the chat
	msg, _, err := wsutil.ReadServerData(conn01)
	assert.Equal(s.T(), "user02 joins the chat", string(msg))
	assert.NoError(s.T(), err)

	// user01 sends message to the chat
	err = wsutil.WriteClientMessage(conn01, ws.OpText, []byte("Ping"))
	assert.NoError(s.T(), err)

	// user02 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(s.T(), "user01: Ping", string(msg))
	assert.NoError(s.T(), err)

	// user02 sends message to user01
	err = wsutil.WriteClientMessage(conn02, ws.OpText, []byte("Pong"))
	assert.NoError(s.T(), err)

	// user01 gets message from user02
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(s.T(), "user02: Pong", string(msg))
	assert.NoError(s.T(), err)

	// user03 joins the chat
	token, err = codeNicknameInToken("user03")
	assert.NoError(s.T(), err)
	conn03, err := getChat(token)
	assert.NoError(s.T(), err)

	// user01 gets message that user03 joined the chat
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(s.T(), "user03 joins the chat", string(msg))
	assert.NoError(s.T(), err)

	// user02 gets message that user03 joined the chat
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(s.T(), "user03 joins the chat", string(msg))
	assert.NoError(s.T(), err)

	// user01 sends message to the chat
	err = wsutil.WriteClientMessage(conn01, ws.OpText, []byte("Ping 2"))
	assert.NoError(s.T(), err)

	// user02 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn02)
	assert.Equal(s.T(), "user01: Ping 2", string(msg))
	assert.NoError(s.T(), err)

	// user02 sends message to user01
	err = wsutil.WriteClientMessage(conn02, ws.OpText, []byte("Pong 2"))
	assert.NoError(s.T(), err)

	// user01 gets message from user02
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(s.T(), "user02: Pong 2", string(msg))
	assert.NoError(s.T(), err)

	// user03 gets message from user01
	msg, _, err = wsutil.ReadServerData(conn03)
	assert.Equal(s.T(), "user01: Ping 2", string(msg))
	assert.NoError(s.T(), err)

	// user03 sends message to user01
	err = wsutil.WriteClientMessage(conn03, ws.OpText, []byte("Pong 2"))
	assert.NoError(s.T(), err)

	// user01 gets message from user03
	msg, _, err = wsutil.ReadServerData(conn01)
	assert.Equal(s.T(), "user03: Pong 2", string(msg))
	assert.NoError(s.T(), err)
}
