package main

import (
	"console-chat/internal/app"
	"console-chat/internal/ports/ginserver"
	"console-chat/internal/ports/wsserver"
	userrepo "console-chat/internal/repo/user_repo"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigFile("configs/config.yml")
	return viper.ReadInConfig()
}

// UserRepoConfig initializes connection to users database
func UserRepoConfig(ctx context.Context, dbURL string) *pgx.Conn {
	// connecting to a database in the loop with delay 1 sec for correct starting in docker container
	for {
		conn, err := pgx.Connect(ctx, dbURL)
		if err != nil { // database haven't initialized in docker container yet
			log.Printf("user_repo connection error: %s\n", err.Error())
			time.Sleep(time.Second)
		} else { // database already initialized
			return conn
		}
	}
}

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("config init error:", err.Error())
	}
	host := viper.GetString("server.ginserver.host")
	port := viper.GetInt("server.ginserver.port")
	tokenKey := []byte(randomdata.Paragraph())

	// configuring userRepo
	userRepoURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("userrepo.username"),
		viper.GetString("userrepo.password"),
		viper.GetString("userrepo.host"),
		viper.GetString("userrepo.port"),
		viper.GetString("userrepo.dbname"),
		viper.GetString("userrepo.sslmode"))

	ctx := context.Background()
	userRepoConn := UserRepoConfig(ctx, userRepoURL)
	defer func() {
		if err := userRepoConn.Close(ctx); err != nil {
			log.Fatal("can't close database connection:", err.Error())
		}
	}()

	ws := wsserver.New(tokenKey)
	app := app.New(userrepo.New(userRepoConn))
	server := ginserver.NewHTTPServer(host, port, ws, app, tokenKey)

	// preparing graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT)

	go func() {
		log.Println("Starting http server")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("can't listen and serve server:", err.Error())
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
