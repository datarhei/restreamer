package rand

// https://www.calhoun.io/creating-random-strings-in-go/

import (
	"math/rand"
	"sync"
	"time"
)

const (
	CharsetLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetNumbers = "1234567890"
	CharsetSymbols = "#@+*%&/<>[]()=?!$.,:;-_"

	CharsetAlphanumeric = CharsetLetters + CharsetNumbers
	CharsetAll          = CharsetLetters + CharsetNumbers + CharsetSymbols
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var lock sync.Mutex

func StringWithCharset(length int, charset string) string {
	lock.Lock()
	defer lock.Unlock()

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func StringLetters(length int) string {
	return StringWithCharset(length, CharsetLetters)
}

func StringNumbers(length int) string {
	return StringWithCharset(length, CharsetNumbers)
}

func StringAlphanumeric(length int) string {
	return StringWithCharset(length, CharsetAlphanumeric)
}

func String(length int) string {
	return StringWithCharset(length, CharsetAll)
}
