package shortuuid

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

type base57 struct {
	// alphabet is the character set to construct the UUID from.
	alphabet alphabet
}

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (b base57) Encode(u uuid.UUID) string {
	var num big.Int
	num.SetString(strings.Replace(u.String(), "-", "", 4), 16)

	// Calculate encoded length.
	length := math.Ceil(math.Log(math.Pow(2, 128)) / math.Log(float64(b.alphabet.Length())))

	return b.numToString(&num, int(length))
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (b base57) Decode(u string) (uuid.UUID, error) {
	str, err := b.stringToNum(u)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(str)
}

// numToString converts a number a string using the given alphabet.
func (b *base57) numToString(number *big.Int, padToLen int) string {
	var (
		out   []rune
		digit *big.Int
	)

	alphaLen := big.NewInt(b.alphabet.Length())

	zero := new(big.Int)
	for number.Cmp(zero) > 0 {
		number, digit = new(big.Int).DivMod(number, alphaLen, new(big.Int))
		out = append(out, b.alphabet.chars[digit.Int64()])
	}

	if padToLen > 0 {
		remainder := math.Max(float64(padToLen-len(out)), 0)
		out = append(out, []rune(strings.Repeat(string(b.alphabet.chars[0]), int(remainder)))...)
	}

	reverse(out)

	return string(out)
}

// stringToNum converts a string a number using the given alphabet.
func (b *base57) stringToNum(s string) (string, error) {
	n := big.NewInt(0)

	for _, char := range s {
		n.Mul(n, big.NewInt(b.alphabet.Length()))

		index, err := b.alphabet.Index(char)
		if err != nil {
			return "", err
		}

		n.Add(n, big.NewInt(index))
	}

	if n.BitLen() > 128 {
		return "", fmt.Errorf("number is out of range (need a 128-bit value)")
	}

	return fmt.Sprintf("%032x", n), nil
}

// reverse reverses a inline.
func reverse(a []rune) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
