package main

import (
	"fmt"
	"jrsa/src"
)

func main() {
	rsa := rsa.NewRsa()
	
	rsa.GenerateKeys()

	msg := "oi"
	pk := rsa.GetPublicKey()
	msgCifrada := rsa.Cifra(msg, pk)
	

	fmt.Println("String ('oi') Cifrada: ", msgCifrada)
	fmt.Println("String decifrada", rsa.Decifra(msgCifrada))
}
