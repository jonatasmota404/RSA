package rsa

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func Get2PrimeNumbers() (uint64, uint64) {
	
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

	primeNumber1, err := strconv.ParseUint(recordNumber[r.Intn(lenRecord)][1], 10, 64)

	if err != nil {
		panic(err)
	}

	primeNumber2, err := strconv.ParseUint(recordNumber[r.Intn(lenRecord)][1], 10, 64)
	if err != nil {
		panic(err)
	}

	return primeNumber1, primeNumber2
}


func GetENumber(Z int) (int) {

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

		E, err := strconv.ParseUint(recordNumber[r.Intn(len)][1], 10, 64)

		if err != nil {
			panic(err)
		}
	
	
		if mdc(Z, int(E)) == 1 && E > 500 && E < uint64(Z) {
			return int(E)
		}
		
	}

}

func mdc(a int, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}
