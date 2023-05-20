package main

import (
	"fmt"
	"math/big"
	"jrsa/src"
)

func main() {
	var stringCifrada []*big.Int
	var stringDecifrada []byte


	stringTeste := "olá mundo"
	stringBytes := []byte(stringTeste)

	fmt.Println("String que será cifrada", stringTeste)

	/*
		Para gerar a sua chave pública, o servidor precisa gerar dois números P e Q.
		Aleatórios, muito grandes e que sejam primos.
	*/
	P := big.NewInt(rsa.GetPrimeNumber())
	Q := big.NewInt(rsa.GetPrimeNumber())

	//Agora calcularemos o N, sendo a multiplicação de P e Q.
	N := big.NewInt(0).Mul(P,Q)

	//Agora será calculado Z que é phi(N) = phi(P) * phi(Q). Ou Z = (P-1) * (Q-1)
	Z := big.NewInt(0).Mul(P.Sub(P, big.NewInt(1)), Q.Sub(Q, big.NewInt(1)))

	//Agora vamos calcular o número E tal que 1 < E < Phi(N)

	E := big.NewInt(rsa.GetENumber(Z.Int64()))

	//A chave pública é composta pelo N e o E.

	D := big.NewInt(rsa.GetDNumber(E.Int64(), Z.Int64())) 

	fmt.Println("String Cifrada")

	for _, v := range stringBytes {
		valorStringEmByte := big.NewInt(int64(v))
		stringCifrada = append(stringCifrada, valorStringEmByte.Exp(valorStringEmByte, E, N))
	}
	fmt.Println(stringCifrada)

	fmt.Println("Sting Decifrada")
	

	for _, v := range stringCifrada {
		numDecifrado := v.Exp(v, D, N)
		stringDecifrada = append(stringDecifrada, numDecifrado.Bytes()...)
	}
	fmt.Println(string(stringDecifrada))
}
