package rsa

import (
	b64 "encoding/base64"
	"encoding/csv"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type RSA struct {
	P *big.Int
	Q *big.Int
	N *big.Int
	Z *big.Int
	E *big.Int
	D *big.Int
}

func NewRsa() *RSA {
	rsa := new(RSA)

	/*
		Para gerar a sua chave pública, o servidor precisa gerar dois números P e Q.
		Aleatórios, muito grandes e que sejam primos.
	*/
	rsa.P = rsa.generatePrimeNumber()
	rsa.Q = rsa.generatePrimeNumber()
	fmt.Println("P = ",rsa.P)
	fmt.Println("Q = ",rsa.Q)

	//Agora calcularemos o N, sendo a multiplicação de P e Q.
	rsa.N = big.NewInt(0).Mul(rsa.P, rsa.Q)
	fmt.Println("N = ",rsa.N)
	//Agora será calculado Z que é phi(N) = phi(P) * phi(Q). Ou Z = (P-1) * (Q-1)
	rsa.Z = big.NewInt(0).Mul(rsa.P.Sub(rsa.P, big.NewInt(1)), rsa.Q.Sub(rsa.Q, big.NewInt(1)))
	fmt.Println("Z = ",rsa.Z)

	//Agora vamos calcular o número E tal que 1 < E < Phi(N)
	
	rsa.E = rsa.generateENumber()
	fmt.Println("E = ",rsa.E)

	rsa.D = rsa.generateDNumber()
	fmt.Println("D = ",rsa.D)
	return rsa
}

func (rsa *RSA) Cifra(msg string, pk string) (string) {
	var stringCifradaSlice []string
	msgByte := []byte(msg)

	pkDecode, _ := b64.StdEncoding.DecodeString(pk)
	pkDecodeSlice := strings.Split(string(pkDecode), "|")
	stringN, _ := strconv.ParseInt(pkDecodeSlice[0], 10, 64)
	stringE, _ := strconv.ParseInt(pkDecodeSlice[1], 10, 64)
	N, E := big.NewInt(stringN), big.NewInt(stringE)

	for _, v := range msgByte {
		valorStringEmByte := big.NewInt(int64(v))
		stringCifradaSlice = append(stringCifradaSlice, valorStringEmByte.Exp(valorStringEmByte, E, N).String())
	}
	stringCifrada := strings.Join(stringCifradaSlice, "|")
	return b64.StdEncoding.EncodeToString([]byte(stringCifrada))
}

func (rsa *RSA) Decifra(msg string) (string) {
	var stringDecifrada []byte

	msgDecode, _ := b64.StdEncoding.DecodeString(msg)
	decodedD , _ := b64.StdEncoding.DecodeString(rsa.getPrivateKey())

	decodedPubKey, _ := b64.StdEncoding.DecodeString(rsa.GetPublicKey())
	decodedPubKeySplited := strings.Split(string(decodedPubKey), "|")
	intN, _ := strconv.ParseInt(decodedPubKeySplited[0], 10, 64)
	N := big.NewInt(intN)
	
	msgDecodeSplited := strings.Split(string(msgDecode), "|")
	intD, _  := strconv.ParseInt(string(decodedD), 10, 64)
	D := big.NewInt(intD)

	for _, v := range msgDecodeSplited {
		numParsedInt, _ := strconv.ParseInt(v, 10, 64)
		numDecifrado := big.NewInt(numParsedInt)
		
		stringDecifrada = append(stringDecifrada, numDecifrado.Exp(numDecifrado, D, N).Bytes()...)
	}

	return string(stringDecifrada)
}

func (rsa *RSA) GetPublicKey() (string) {
	pk, err := os.ReadFile("pk")
	if err != nil {
		panic(err)	
	}
	pubKeyDecode, _ := b64.StdEncoding.DecodeString(string(pk))
	pubKeySlice := strings.Split(string(pubKeyDecode), "|")
	publicKey := pubKeySlice[0] + "|" + pubKeySlice[1]
	return b64.StdEncoding.EncodeToString([]byte(publicKey))
}

func (rsa *RSA) GenerateKeys() {
	f := MakeFile("pk")
	defer CloseFile(f)
	keysString := rsa.N.String() + "|" + rsa.E.String() + "|" + rsa.D.String()
	f.WriteString(b64.StdEncoding.EncodeToString([]byte(keysString)))
}

func (rsa *RSA) generatePrimeNumber() (*big.Int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	file, err := os.Open("primes.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	recordNumber, err := csv.NewReader(file).ReadAll()

	if err != nil {
		panic(err)
	}

	lenRecord := len(recordNumber) + 1

	primeNumber, err := strconv.ParseInt(recordNumber[r.Intn(lenRecord)][1], 10, 64)

	if err != nil {
		panic(err)
	}

	return big.NewInt(primeNumber)
}

func (rsa *RSA) generateENumber() (*big.Int) {
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	x := rsa.Z.Int64()

	for {
		E := big.NewInt(int64(r.Intn(int(x))))

		if AreCoprime(rsa.Z, E) {
			return E
		}
	}
}

func (rsa *RSA) getPrivateKey() (string) {
	pk, err := os.ReadFile("pk")
	if err != nil {
		panic(err)	
	}

	priKeyDecode, _ := b64.StdEncoding.DecodeString(string(pk))
	priKeySlice := strings.Split(string(priKeyDecode), "|")
	privateKey := priKeySlice[2]
	return b64.StdEncoding.EncodeToString([]byte(privateKey))
}

func (rsa *RSA) generateDNumber() (*big.Int){
	/*
	Usar algoritmos Euclidianos Estendidos que pegam dois inteiros 'a' e 'b', então encontre seu mdc, e também encontre 'x' e 'y' tal que 
	ax + by = mdc(a, b)

	https://www.geeksforgeeks.org/multiplicative-inverse-under-modulo-m/
	*/
	D := big.NewInt(0)
	return D.ModInverse(rsa.E, rsa.Z)
}