package main

import (
	"console-chat/internal/ports/ginserver"
	"console-chat/internal/ports/wsserver"
	"errors"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigFile("config.yml")
	return viper.ReadInConfig()
}

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("config init error:", err.Error())
	}
	host := viper.GetString("server.ginserver.host")
	port := viper.GetInt("server.ginserver.port")

	ws := wsserver.NewWsServer()
	server := ginserver.NewHTTPServer(host, port, ws)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("can't listen and serve server: %s", err.Error())
	}
}
