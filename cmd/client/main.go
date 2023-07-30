package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"golang.org/x/crypto/ssh/terminal"
)

const httpUrl = "http://localhost:8080/console-chat/users"
const wsUrl = "ws://localhost:8080/console-chat/chat"

type registerResponse struct {
	Data struct {
		Nickname       string `json:"nickname"`
		HashedPassword string `json:"hashed_password"`
	} `json:"data"`
	Error string `json:"error"`
}

// RegisterNewUser gets new user nickname & password from stdin and makes http request to register new user
func RegisterNewUser() {
	for {
		// getting new user nickname
		fmt.Print("Enter new user's nickname: ")
		reader := bufio.NewReader(os.Stdin)
		nickname, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("nickname input error:", err.Error())
		}
		nickname = nickname[:len(nickname)-1]

		// getting new user password
		fmt.Print("Enter new user's password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("password input error:", err.Error())
		}
		fmt.Println()

		// check if the message to http server is valid
		jsonStr := fmt.Sprintf(`{
			"nickname": "%s",
			"password": "%s"
		}`, nickname, string(password))
		if !json.Valid([]byte(jsonStr)) {
			fmt.Println("Don't use symbols like \", ` or '. Please try again.")
			continue
		}

		// making request to the http server and getting the response
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, httpUrl, strings.NewReader(jsonStr))
		if err != nil {
			log.Fatal("request error:", err.Error())
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			log.Fatal("request execution error:", err.Error())
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Fatal("response body closure error:", err.Error())
			}
		}()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("response body read error:", err.Error())
		}
		var regResp registerResponse
		if err := json.Unmarshal(body, &regResp); err != nil {
			log.Fatal("response body unmarshalling error:", err.Error())
		}

		// checking result
		if regResp.Error == "user has invalid password" || regResp.Error == "user has invalid nickname" {
			fmt.Println("Nickname should have length between 4 and 25 including borders, contain only latin letters, digits or _ and don't contain obscenities")
			fmt.Println("Password should have length between 6 and 50 including borders, contain at least 1 latin letter, 1 digit and 1 special symbol")
			continue
		} else if regResp.Error == "user with required nickname already exists" {
			fmt.Println("User with nickname", nickname, "already exists")
			continue
		} else if regResp.Error != "" {
			fmt.Println("Please try again later")
			return
		} else {
			fmt.Println("Registration successfully completed")
			return
		}
	}
}

func IsValidUrl(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

type signInResponse struct {
	Data struct {
		TokenString string `json:"token_string"`
	} `json:"data"`
	Error string `json:"error"`
}

func SignIn() string {
	for {
		// getting user nickname
		fmt.Print("Enter your nickname: ")
		reader := bufio.NewReader(os.Stdin)
		nickname, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("nickname input error:", err.Error())
		}
		nickname = nickname[:len(nickname)-1]

		// getting user password
		fmt.Print("Enter your password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal("password input error:", err.Error())
		}
		fmt.Println()

		// check if the message to http server is valid
		jsonStr := fmt.Sprintf(`{
			"password": "%s"
		}`, string(password))
		getUrl := httpUrl + "/" + nickname
		if !(json.Valid([]byte(jsonStr)) && IsValidUrl(getUrl)) {
			fmt.Println("Don't use symbols like \", ` or '. Please try again.")
			continue
		}

		// making request to the http server and getting the response
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, getUrl, strings.NewReader(jsonStr))
		if err != nil {
			log.Fatal("request error:", err.Error())
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			log.Fatal("request execution error:", err.Error())
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Fatal("response body closure error:", err.Error())
			}
		}()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("response body read error:", err.Error())
		}
		var signResp signInResponse
		if err := json.Unmarshal(body, &signResp); err != nil {
			log.Fatal("response body unmarshalling error:", err.Error())
		}

		// checking result
		switch signResp.Error {
		case "could not find required user":
			fmt.Println("User with nickname", nickname, "doesn't exist")
			continue
		case "wrong password of required user":
			fmt.Println("Wrong password of user with nickname", nickname)
			continue
		case "":
			fmt.Println("Successfully signed in")
			return signResp.Data.TokenString
		}
	}
}

func main() {
	reg := flag.Bool("reg", false, "Flag to register new user")
	sign := flag.Bool("sign", false, "Flag to sign in and join the chat")
	flag.Parse()

	if *reg == *sign {
		fmt.Println("This is console-chat client. Run this program with \"-reg\" flag to register new user or \"-sign\" flag to sign in and join the chat")
	} else if *reg { // registration of the new user
		RegisterNewUser()
	} else { // signing in and connecting to the chat
		token := SignIn()

		// connecting to websocket server
		conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), wsUrl)
		if err != nil {
			log.Fatalf("can't dial connection: %s", err.Error())
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Fatal("connection closure error:", err.Error())
			}
		}()

		// sending token to authorize
		if err := wsutil.WriteClientMessage(conn, ws.OpText, []byte(token)); err != nil {
			log.Fatal("can't wtite client token:", err.Error())
		} else {
			fmt.Println("Successfully connected to chat. Start writing messages!")
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
}
