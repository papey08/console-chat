package ginserver

import (
	"console-chat/internal/app"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(host string, port int, app app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("console-chat")
	AppRouter(api, app)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}
