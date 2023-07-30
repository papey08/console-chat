package ginserver

import (
	"console-chat/internal/app"
	"console-chat/internal/ports/wsserver"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(host string, port int, ws *wsserver.WsServer, app app.App, tokenKey []byte) *http.Server {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	api := router.Group("console-chat")
	AppRouter(api, ws, app, tokenKey)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}
