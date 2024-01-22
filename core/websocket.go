package core

import (
	"ShopService/helpers"
	"ShopService/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

type Chat struct {
	ID          string
	Connections map[*Connection]bool //to be used to hold connections
	Messages    []Message
	Users       map[models.User]bool
}

type Message struct {
	Sender  *Connection
	Message string
}

type Connection struct {
	WSConn *websocket.Conn
	User   *models.User
}

func (client *Connection) broadcast(msg Message, chat *Chat) {

	for conn := range chat.Connections {
		if err := conn.WSConn.WriteJSON(msg); err != nil {
			delete(chat.Connections, conn)
		}
	}
}

func (client *Connection) echo(msg EchoPayload) {
	if err := client.WSConn.WriteJSON(msg); err != nil {
		panic("An error occurred while sending ws")
	}
	payload, _ := json.Marshal(msg)
	println("Message: ", string(payload))
}

func RunWebsocketServer() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/ws/{token}/chat/{chatID}", chatHandler)
	muxRouter.HandleFunc("/ws/{token}/chats", chatsHandler)
	muxRouter.HandleFunc("/ws/{token}/echo", echoHandler)
	http.Handle("/", muxRouter)
	fmt.Println("WebSocket server is running on port 8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return
	}

}

var chatRooms = make(map[string]*Chat)

func getUserConnection(w http.ResponseWriter, r *http.Request) *Connection {
	paths := mux.Vars(r)
	upgrader := websocket.Upgrader{}
	token := paths["token"]
	claims, err := ValidateAccessToken(token)
	if err != nil {
		panic(err.Error())
	}
	wsConn, _ := upgrader.Upgrade(w, r, nil)
	return &Connection{wsConn, &claims.User}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	paths := mux.Vars(r)
	chatID := paths["chat_id"]
	client := getUserConnection(w, r)
	chat := chatRooms[chatID]
	if !helpers.ContainsItem(chat.Users, client.User) {
		panic("You are not a member of this chat")
	}
	chat.Connections[client] = true
	go func() {
		for {
			msg := Message{}
			err := client.WSConn.ReadJSON(&msg)
			if err != nil {
				return
			}

			chat.Messages = append(chat.Messages)

			client.broadcast(msg, chat)
		}
	}()

}

type EchoPayload struct {
	User  models.User
	Data  string
	Event string
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	client := getUserConnection(w, r)
	go func() {
		for {
			msg := EchoPayload{}
			err := client.WSConn.ReadJSON(&msg)
			if err != nil {
				return
			}
			client.echo(msg)
		}
	}()
}

func chatsHandler(w http.ResponseWriter, r *http.Request) {
	client := getUserConnection(w, r)
	go func() {
		for {
			msg := EchoPayload{}
			err := client.WSConn.ReadJSON(&msg)
			if err != nil {
				return
			}
			client.echo(msg)
		}
	}()
}

// func sendWs()
func connectToWebSocket(url string) (*websocket.Conn, error) {
	// Connect to the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to WebSocket: %v", err)
	}

	return conn, nil
}

type Payload struct {
	User  models.User
	Data  interface{}
	Event string
}

func sendDataToWebSocket(conn *websocket.Conn, data Payload) error {
	// Send data to the WebSocket
	payload, _ := json.Marshal(data)
	err := conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return fmt.Errorf("error sending data to WebSocket: %v", err)
	}

	return nil
}

func SendToWs(route string, payload Payload) {
	token, err0 := GenerateJWT(payload.User, false)
	if err0 != nil {
		panic(err0.Error())
	}
	wsUrl := "ws://localhost:8081/ws/"
	conn, err := connectToWebSocket(fmt.Sprintf("%s%s/%s", wsUrl, token, route))
	if err != nil {
		panic(err.Error())
	}
	err2 := sendDataToWebSocket(conn, payload)
	if err2 != nil {
		panic(err2.Error())
	}
}
