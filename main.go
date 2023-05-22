package main

import (
	"fmt"
	"jrsa/src/rsa"
	"jrsa/src/chat"

)

func main() {
	rsa := rsa.NewRsa()

	msg := "Olá mundo"
	pk := rsa.GetPublicKey()
	msgCifrada := rsa.Cifra(msg, pk)
	

	fmt.Println("String ('"+msg+"') Cifrada: ", msgCifrada)
	fmt.Println("String decifrada", rsa.Decifra(msgCifrada))
	chat.Serve()
}
