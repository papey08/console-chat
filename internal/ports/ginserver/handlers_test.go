package ginserver

import (
	"bytes"
	"console-chat/internal/model"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

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
