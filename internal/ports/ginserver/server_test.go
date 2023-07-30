package ginserver

import (
	mocks "console-chat/internal/ports/ginserver/app_mocks"
	"console-chat/internal/ports/wsserver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ginServerTestSuite struct {
	suite.Suite
	app     *mocks.App
	client  *http.Client
	server  *http.Server
	baseURL string
}

func ginServerTestSuiteInit(s *ginServerTestSuite) {
	s.app = new(mocks.App)

	tokenKey := []byte("abcd")
	ws := wsserver.New(tokenKey)
	s.server = NewHTTPServer("localhost", 8081, ws, s.app, tokenKey)
	testServer := httptest.NewServer(s.server.Handler)
	s.client = testServer.Client()
	s.baseURL = testServer.URL
}

func (s *ginServerTestSuite) SetupSuite() {
	ginServerTestSuiteInit(s)
}

func (s *ginServerTestSuite) TearDownSuite() {
	_ = s.server.Close()
}

func (s *ginServerTestSuite) getResponse(req *http.Request, out any) (int, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("unexpected error: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("unable to read response: %w", err)
	}
	_ = json.Unmarshal(respBody, out)
	return resp.StatusCode, nil
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ginServerTestSuite))
}
