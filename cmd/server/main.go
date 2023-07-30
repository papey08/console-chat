package main

import (
	"console-chat/internal/app"
	"console-chat/internal/ports/ginserver"
	"console-chat/internal/ports/wsserver"
	userrepo "console-chat/internal/repo/user_repo"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pallinder/go-randomdata"
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
	tokenKey := []byte(randomdata.Paragraph())

	ws := wsserver.NewWsServer(tokenKey)
	app := app.New(userrepo.New())
	server := ginserver.NewHTTPServer(host, port, ws, app, tokenKey)

	// preparing graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT)

	go func() {
		log.Println("Starting http server")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("can't listen and serve server: %s", err.Error())
		}
	}()

	<-osSignals

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server graceful shutdown failed:", err.Error())
	}
	log.Println("Server was gracefully stopped")
}
