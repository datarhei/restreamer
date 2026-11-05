package shortuuid

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultAlphabet is the default alphabet used.
const DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type alphabet struct {
	chars [57]rune
	len   int64
}

// Remove duplicates and sort it to ensure reproducability.
func newAlphabet(s string) alphabet {
	abc := dedupe(strings.Split(s, ""))

	if len(abc) != 57 {
		panic("encoding alphabet is not 57-bytes long")
	}

	sort.Strings(abc)
	a := alphabet{
		len: int64(len(abc)),
	}

	for i, char := range strings.Join(abc, "") {
		a.chars[i] = char
	}

	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}

// Index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (a *alphabet) Index(t rune) (int64, error) {
	for i, char := range a.chars {
		if char == t {
			return int64(i), nil
		}
	}
	return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
}

// dudupe removes duplicate characters from s.
func dedupe(s []string) []string {
	var out []string
	m := make(map[string]bool)

	for _, char := range s {
		if _, ok := m[char]; !ok {
			m[char] = true
			out = append(out, char)
		}
	}

	return out
}
