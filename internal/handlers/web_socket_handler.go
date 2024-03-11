package handlers

import (
	"fmt"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	privateConns     = make(map[string]*websocket.Conn)
	privateBroadcast = make(chan client.PrivateMessage)

	groupConns = make(map[string][]*websocket.Conn)
)

func HandlePrivateConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		var message client.PrivateMessage
		if err := conn.ReadJSON(&message); err != nil {
			fmt.Println(err)
			break
		}

		if _, ok := privateConns[message.Sender]; !ok {
			privateConns[message.Sender] = conn
		}

		privateBroadcast <- message
	}
}

func HandlePrivateMessages() {
	for {
		message := <-privateBroadcast
		if conn, ok := privateConns[message.Receiver]; ok {
			if err := conn.WriteJSON(&message); err != nil {
				fmt.Println(err)
				conn.Close()
				delete(privateConns, message.Receiver)
			}
		}
	}
}

func HandleGroupMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		var message client.GroupMessage
		if err := conn.ReadJSON(&message); err != nil {
			fmt.Println(err)
			break
		}

		if conns, ok := groupConns[message.Group]; ok {
			var isPresent bool
			for _, groupConn := range conns {
				if groupConn == conn {
					isPresent = true
					break
				}
			}
			if !isPresent {
				groupConns[message.Group] = append(groupConns[message.Group], conn)
			}
		} else {
			groupConns[message.Group] = append(groupConns[message.Group], conn)
		}

		if conns, ok := groupConns[message.Group]; ok {
			for _, groupConn := range conns {
				if groupConn != conn {
					if err := groupConn.WriteJSON(&message); err != nil {
						fmt.Println(err)
						groupConn.Close()
						delete(privateConns, message.Group)
					}
				}
			}
		}
	}
}

// // Upgrade the HTTP protocol to Web Socket protocol
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// var p2pConnections = make(map[string]*websocket.Conn)

// // Handler to process bidirectional requests on private chats
// func HandlePrivateChat(w http.ResponseWriter, r *http.Request) {
// 	logger := configs.NewLogger("handlers")

// 	// Check if token is valid
// 	claims, err := handleAuthorization(w, r)
// 	if err != nil {
// 		return
// 	}

// 	// Update connection to web socket
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		logger.Debugf("Failed to upgrade to web socket: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	p2pConnections[claims.Username] = conn
// 	defer delete(p2pConnections, claims.Username)

// 	for {
// 		var body client.P2PConn
// 		if err := conn.ReadJSON(&body); err != nil {
// 			logger.Debug(err)
// 			break
// 		}

// 		channel := saveMessage(&body)

// 		if receiverConn, ok := p2pConnections[body.Receiver]; ok {
// 			err = receiverConn.WriteJSON(&body)
// 			if err != nil {
// 				fmt.Println(err)
// 				if err = conn.WriteJSON(client.WebSocketError{Err: err.Error()}); err != nil {
// 					fmt.Println(err)
// 				}
// 				break
// 			}
// 			if err = conn.WriteJSON(&body); err != nil {
// 				fmt.Println(err)
// 			}
// 		} else {
// 			result := <-channel
// 			if result.Err != nil {
// 				// record not found
// 				fmt.Println(result.Err)
// 				if err = conn.WriteJSON(client.WebSocketError{Err: result.Err.Error()}); err != nil {
// 					fmt.Println(err)
// 				}
// 				break
// 			}
// 			go changeMessageStatus(result.Value)
// 			if err = conn.WriteJSON(client.WebSocketError{Err: "User is not online, message will be send"}); err != nil {
// 				fmt.Println(err)
// 			}
// 			break
// 		}
// 	}
// }

// // Handler to process bidirectional requests on group chats
// func HandleGroupChat(w http.ResponseWriter, r *http.Request) {
// 	logger := configs.NewLogger("handlers")

// 	// Check if token is valid
// 	_, err := handleAuthorization(w, r)
// 	if err != nil {
// 		return
// 	}

// 	// Update connection to web socket
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		logger.Debugf("Failed to upgrade to web socket: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			logger.Debug(err)
// 			break
// 		}

// 		err = conn.WriteMessage(messageType, p)
// 		if err != nil {
// 			logger.Debug(err)
// 			break
// 		}
// 	}
// }

// type ChannelResult struct {
// 	Value *models.PrivateMessage
// 	Err   error
// }

// func saveMessage(body *client.P2PConn) <-chan ChannelResult {
// 	channel := make(chan ChannelResult)
// 	go func() {
// 		sender, err := repositories.FindUserByName(body.Sender)
// 		if err != nil {
// 			channel <- ChannelResult{nil, err}
// 			close(channel)
// 			return
// 		}

// 		receiver, err := repositories.FindUserByName(body.Receiver)
// 		if err != nil {
// 			channel <- ChannelResult{nil, err}
// 			close(channel)
// 			return
// 		}

// 		message, err := repositories.SavePrivateMessage(&models.PrivateMessage{Sender: *sender, Receiver: *receiver, Data: body.Data, IsPending: false})
// 		if err != nil {
// 			channel <- ChannelResult{nil, err}
// 			close(channel)
// 			return
// 		}

// 		channel <- ChannelResult{Value: message, Err: nil}
// 		close(channel)
// 	}()
// 	return channel
// }

// func changeMessageStatus(message *models.PrivateMessage) {
// 	message.IsPending = true
// 	repositories.UpdatePendingSituation(message)
// }
