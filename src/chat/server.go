package chat

import (
    "fmt"
    "log"
    "net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}


func chatHandler(w http.ResponseWriter, r *http.Request) {

    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }

	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

    for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		
		fmt.Println(string(msg))
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}


func Serve() {
	fileServer := http.FileServer(http.Dir("./public"))
    http.Handle("/", fileServer)
    http.HandleFunc("/chat", chatHandler)



    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

