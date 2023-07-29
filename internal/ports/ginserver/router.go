package ginserver

import (
	"console-chat/internal/ports/wsserver"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, ws *wsserver.WsServer) {
	r.GET("/chat", gin.WrapF(ws.Chat))
}
