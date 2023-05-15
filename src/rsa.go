package rsa

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func Get2PrimeNumbers() (int64, int64) {
	
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

	primeNumber1, err := strconv.ParseInt(recordNumber[r.Intn(lenRecord)][1], 10, 64)

	if err != nil {
		panic(err)
	}

	primeNumber2, err := strconv.ParseInt(recordNumber[r.Intn(lenRecord)][1], 10, 64)
	if err != nil {
		panic(err)
	}

	return primeNumber1, primeNumber2
}


func GetENumber(Z int64) (int64) {

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

	len := len(recordNumber) + 1

	for {

		E, err := strconv.ParseInt(recordNumber[r.Intn(len)][1], 10, 64)

		if err != nil {
			panic(err)
		}
	
	
		if mdc(Z, E) == 1 && E > 500 && E < Z {
			return E
		}

	}
}

func GetDNumber(E, Z int64) (int64){
	return modInverse(E, Z)
}

func mdc(a, b int64) int64 {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

/*
Usar algoritmos Euclidianos Estendidos que pegam dois inteiros 'a' e 'b', então encontre seu mdc, e também encontre 'x' e 'y' tal que 
ax + by = mdc(a, b)

https://www.geeksforgeeks.org/multiplicative-inverse-under-modulo-m/
*/

func modInverse(A, M int64) int64 {
    var y int64 = 0
    var x int64 = 1
	m0 := M

    if M == 1 {
        return 0
    }

    for A > 1 {

        // q is quotient
        q := A / M

        t := M

        // m is remainder now, process
        // same as Euclid's algo
        M = A % M
        A = t
        t = y

        // Update x and y
        y = x - q*y
        x = t
    }

    // Make x positive
    if x < 0 {
        x += m0
    }

    return x
}