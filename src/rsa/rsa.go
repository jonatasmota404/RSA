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
//Estrutura para criação de um objeto rsa.
//(Obs em golang as letras minúsculas são atributos privados, já as maiúsculas são publicos)
type RSA struct {
	p *big.Int
	q *big.Int
	N *big.Int
	Z *big.Int
	E *big.Int
	d *big.Int
}

//Método construtor do objeto rsa
func NewRsa() *RSA {
	rsa := new(RSA)

	//verifica se as chaves já foram geradas
	if !FileExists("pk") {
		//gera as chaves e retorna o obj rsa
		rsa.GenerateKeys()
		return rsa
	}

	//atribui os parâmetros da chave pública e privada as suas respectivas variáveis
	rsa.N, rsa.E, rsa.d = rsa.getPkParameters()
	return rsa
}

//cifra a mensagem de acordo com a chave pública
//obs : a cifra tem que ser de acordo com a chave pública compartilhada pelo destinatário
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

//decifra a mensagem de acordo com nossa chave privada
func (rsa *RSA) Decifra(msg string) (string) {
	var stringDecifrada []byte

	msgDecode, _ := b64.StdEncoding.DecodeString(msg)
	msgDecodeSplited := strings.Split(string(msgDecode), "|")

	
	for _, v := range msgDecodeSplited {
		numParsedInt, _ := strconv.ParseInt(v, 10, 64)
		numDecifrado := big.NewInt(numParsedInt)
		
		stringDecifrada = append(stringDecifrada, numDecifrado.Exp(numDecifrado, rsa.d, rsa.N).Bytes()...)
	}

	return string(stringDecifrada)
}

//fornece nossa chave pública
func (rsa *RSA) GetPublicKey() (string) {
	N, E, _ := rsa.getPkParameters()
	publicKey := N.String() + "|" + E.String()
	return b64.StdEncoding.EncodeToString([]byte(publicKey))
}

//gera as chaves
func (rsa *RSA) GenerateKeys() {
	f := MakeFile("pk")
	defer CloseFile(f)

	rsa.p = rsa.generatePrimeNumber()
	rsa.q = rsa.generatePrimeNumber()

	fmt.Println("P = ",rsa.p)
	fmt.Println("Q = ",rsa.q)

	//Agora calcularemos o N, sendo a multiplicação de P e Q.
	rsa.N = big.NewInt(0).Mul(rsa.p, rsa.q)
	fmt.Println("N = ",rsa.N)
	//Agora será calculado Z que é phi(N) = phi(P) * phi(Q). Ou Z = (P-1) * (Q-1)
	rsa.Z = big.NewInt(0).Mul(rsa.p.Sub(rsa.p, big.NewInt(1)), rsa.q.Sub(rsa.q, big.NewInt(1)))
	fmt.Println("Z = ",rsa.Z)

	//Agora vamos calcular o número E tal que 1 < E < Phi(N)
	
	rsa.E = rsa.generateENumber()
	fmt.Println("E = ",rsa.E)

	rsa.d = rsa.generateDNumber()
	fmt.Println("D = ",rsa.d)

	keysString := rsa.N.String() + "|" + rsa.E.String() + "|" + rsa.d.String()
	f.WriteString(b64.StdEncoding.EncodeToString([]byte(keysString)))
}

//gera um número primo aleatório
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

//gera a variável E
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

//gera a variável D
func (rsa *RSA) generateDNumber() (*big.Int){
	/*
	Usar algoritmos Euclidianos Estendidos que pegam dois inteiros 'a' e 'b', então encontre seu mdc, e também encontre 'x' e 'y' tal que 
	ax + by = mdc(a, b)

	https://www.geeksforgeeks.org/multiplicative-inverse-under-modulo-m/
	*/
	D := big.NewInt(0)
	return D.ModInverse(rsa.E, rsa.Z)
}

//fornece os parâmetros da chave
func (rsa *RSA) getPkParameters() (*big.Int,*big.Int,*big.Int) {
	pk, err := os.ReadFile("pk")
	if err != nil {
		panic(err)	
	}
	pubKeyDecode, _ := b64.StdEncoding.DecodeString(string(pk))
	pubKeySlice := strings.Split(string(pubKeyDecode), "|")
	N, _ := strconv.ParseInt(pubKeySlice[0], 10, 64)
	E, _ := strconv.ParseInt(pubKeySlice[1], 10, 64)
	D, _ := strconv.ParseInt(pubKeySlice[2], 10, 64)

	return big.NewInt(N), big.NewInt(E), big.NewInt(D)
}