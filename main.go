package main

import (
	"fmt"
	"jrsa/src"
)

func main() {
	stringTeste := "olá mundo"
	stringBytes := []byte(stringTeste)
	fmt.Println(stringBytes)

	
	/*
		Para gerar a sua chave pública, o servidor precisa gerar dois números P e Q.
		Aleatórios, muito grandes e que sejam primos.
	*/
	P, Q := rsa.Get2PrimeNumbers()

	//Agora calcularemos o N, sendo a multiplicação de P e Q.
	N := P * Q

	//Agora será calculado Z que é phi(N) = phi(P) * phi(Q). Ou Z = (P-1) * (Q-1)
	Z := (P - 1) * (Q - 1)

	//Agora vamos calcular o número E tal que 1 < E < Phi(N)

	E := rsa.GetENumber(Z)

	//A chave pública é composta pelo N e o E.

	D := rsa.GetDNumber(E, Z)


	var stringCifrada []int64
	var stringDecifrada []byte


	for _, b := range stringBytes {
		strInt := int64(b)
		stringCifrada = append(stringCifrada, E**&strInt)
	}

	fmt.Println("String Cifrada")	

	fmt.Println(stringCifrada)	


	fmt.Println("String deficafrada")

	for _, s := range stringCifrada {
		stringDecifrada = append(stringDecifrada, uint8((D**&s)%N))
	}

	fmt.Println(stringDecifrada)
}
