package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

const url = "ws://localhost:8080/console-chat/chat"

type userPresenter struct {
	Nickname string `json:"nickname"`
}

func main() {
	// connecting to websocket server
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		log.Fatalf("can't dial connection: %s", err.Error())
	}

	// getting info about the user
	fmt.Print("Enter your name: ")
	var usr userPresenter
	fmt.Scan(&usr.Nickname)

	// sending user info to the server
	usrData, err := json.Marshal(usr)
	if err != nil {
		log.Fatal("user nickname marshalling error:", err.Error())
	}
	if err := wsutil.WriteClientMessage(conn, ws.OpText, usrData); err != nil {
		log.Fatal("can't send user info to server:", err.Error())
	}

	// reading messages from the server
	go func() {
		for {
			data, _, err := wsutil.ReadServerData(conn)
			if err != nil && err != io.EOF {
				log.Fatal("can't read server data:", err.Error())
			} else if err == io.EOF {
				log.Println("server stopped")
				os.Exit(0)
			}

			fmt.Println(string(data))
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	// sending messages from the user
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("can't read string from stdin: %s", err.Error())
		}

		text = text[:len(text)-1]

		if err := wsutil.WriteClientMessage(conn, ws.OpText, []byte(text)); err != nil {
			log.Fatal("can't wtite client message:", err.Error())
		}

		time.Sleep(100 * time.Millisecond) // delay between sending messages
	}
}
