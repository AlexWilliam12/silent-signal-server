package main

import (
	"fmt"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/gorilla/websocket"
)

func main() {

	var token string
	fmt.Print("Your token: ")
	if _, err := fmt.Scanln(&token); err != nil {
		panic(err)
	}

	var sender string
	fmt.Print("Your name: ")
	if _, err := fmt.Scanln(&sender); err != nil {
		panic(err)
	}

	var receiver string
	fmt.Print("Receiver name: ")
	if _, err := fmt.Scanln(&receiver); err != nil {
		panic(err)
	}

	url := "ws://localhost:8080/chat/private"
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		for {
			var response client.PrivateMessage
			err := conn.ReadJSON(&response)
			if err != nil {
				panic(err)
			}

			fmt.Println(response.Message)
		}
	}()

	for {
		var input string
		_, err := fmt.Scanln(&input)

		if err != nil {
			panic(err)
		}

		if input == "/exit" {
			break
		}

		err = conn.WriteJSON(client.PrivateMessage{
			Sender:   sender,
			Receiver: receiver,
			Message: client.Message{
				Type:    "text",
				Content: input,
			}})
		if err != nil {
			panic(err)
		}

	}
}
