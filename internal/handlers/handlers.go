package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var clients = make(map[WebSocketConnection]string)
var wsChan = make(chan WsPayload)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

func Home(w http.ResponseWriter, r *http.Request) {
	// Fetch the template
	view, err := views.GetTemplate("home.jet")
	if err != nil {
		log.Println("Error getting template:", err)
		return
	}

	err = view.Execute(w, nil, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	log.Println("Client connected to endpoint!")
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`
	conn.WriteJSON(response)
	go ListenForWs(&conn)

}

func ListenForWs(conn *WebSocketConnection) {
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// Client disconnected
			payload.Action = "left"
			payload.Conn = *conn
			wsChan <- payload
			break
		}
		payload.Conn = *conn
		wsChan <- payload
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err:", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {

		e := <-wsChan

		switch e.Action {
		case "broadcast":
			response.Action = "broadcast"
			response.Message = "<strong>" + e.Username + "</strong>: " + e.Message
			broadcastToAll(response)
		case "username":
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			delete(clients, e.Conn)
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		}
	}
}

// getUserList helps give back a slice of uniquely connected usernames
func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	return userList
}
