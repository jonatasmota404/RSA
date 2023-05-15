package main

import (
	"fmt"
	"jrsa/src"
)

func main() {
	/*
	Para gerar a sua chave pública, o servidor precisa gerar dois números P e Q.
	Aleatórios, muito grandes e que sejam primos. 
	*/
	P, Q :=  rsa.Get2PrimeNumbers()

	//Agora calcularemos o N, sendo a multiplicação de P e Q.
	N := P * Q

	//Agora será calculado Z que é phi(N) = phi(P) * phi(Q). Ou Z = (P-1) * (Q-1)
	Z := (P - 1) * (Q - 1)

	//Agora vamos calcular o número E tal que 1 < E < Phi(N)

	E := rsa.GetENumber(Z)

	//A chave pública é composta pelo N e o E.

	D := rsa.GetDNumber(E, Z)

	fmt.Println(P, Q, N, Z, E, D)
}