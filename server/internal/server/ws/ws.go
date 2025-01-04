package ws

import (
    "net/http"
	"fmt"
    "github.com/gorilla/websocket"
    "github.com/golang-jwt/jwt/v5"
	"sync"
)
//
// A map of connections between the server and clients

var WSConnMutex = &sync.Mutex{}
var WSConnections = make(map[string]*websocket.Conn)  // Map of UUID to WebSocket connection

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Adjust this to a more secure setting as needed
    },
}


func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket connection requested")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
    if !ok {
        http.Error(w, "Could not get claims from context", http.StatusInternalServerError)
        return
    }
	uuid := claims["sub"].(string)

	WSConnMutex.Lock()
    WSConnections[uuid] = ws
    WSConnMutex.Unlock()

	fmt.Println("Websocket connection opened", uuid)
    defer func() {
        WSConnMutex.Lock()
        delete(WSConnections, uuid)
        WSConnMutex.Unlock()
		fmt.Println("Websocket connection closed", uuid)
    }()

	select {
    case <-r.Context().Done():
        fmt.Println("WebSocket closed by client or context timeout")
    }
	// ListenForDBINserts will send a message to the ws if there is a new db insert by uuid
}



// // Use claims
// fmt.Println("Access granted. User ID: ", claims["sub"])

// fmt.Println("Websocket connection opened")
// word := "sadf"
// time.Sleep(10 * time.Second)  // Example delay


// // Assume you have some logic here to wait for a new database entry
// msg := word
// fmt.Println("Sending message", msg)
// if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
// 	log.Println(err)
// }