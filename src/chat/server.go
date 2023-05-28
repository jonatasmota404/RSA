package chat

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"jrsa/src/rsa"
	"log"
	"net/http"
	"strings"

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

		//rsa := rsa.NewRsa()
		//msgDecoded := rsa.Decifra(string(msg))
		err = conn.WriteMessage(mt, []byte(msg))
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}

func chatHandlerPk(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/chat/pk" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	rsa := rsa.NewRsa()
	pk := rsa.GetPublicKey()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.Header().Set("Access-Control-Allow-Origin", "*")          // normal header
	w.Header().Set("Access-Control-Allow-Credentials", "true")  // normal header
	w.Header().Set("Access-Control-Allow-Methods", "GET")       // normal header
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "{\"pk\": \"%s\"}", pk)
}

func chatHandlerCifra(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	if r.URL.Path != "/chat/cifra" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain;") // normal header
	w.Header().Set("Access-Control-Allow-Origin", "*")          // normal header
	w.Header().Set("Access-Control-Allow-Credentials", "true")  // normal header
	w.Header().Set("Access-Control-Allow-Methods", "POST")       // normal header
	w.WriteHeader(http.StatusOK)

	msgDecoded, _ := b64.StdEncoding.DecodeString(string(body))
	msgDecodeSplited := strings.Split(string(msgDecoded), "|")
		
	rsa := rsa.NewRsa()
	msgCifrada := rsa.Cifra(msgDecodeSplited[0], msgDecodeSplited[1])
	fmt.Fprintf(w, "{\"msg\": \"%s\"}", msgCifrada)	
}

func chatHandlerDecifra(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	if r.URL.Path != "/chat/decifra" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain;") // normal header
	w.Header().Set("Access-Control-Allow-Origin", "*")          // normal header
	w.Header().Set("Access-Control-Allow-Credentials", "true")  // normal header
	w.Header().Set("Access-Control-Allow-Methods", "POST")       // normal header
	w.WriteHeader(http.StatusOK)

	rsa := rsa.NewRsa()
	msgDecoded := rsa.Decifra(string(body))

	fmt.Fprintf(w, "{\"msg\": \"%s\"}", msgDecoded)	
}

func Serve() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/", fileServer)

	http.HandleFunc("/chat/cifra", chatHandlerCifra)
	http.HandleFunc("/chat/decifra", chatHandlerDecifra)

	http.HandleFunc("/chat/pk", chatHandlerPk)
	http.HandleFunc("/chat", chatHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
