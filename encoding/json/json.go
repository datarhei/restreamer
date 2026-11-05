// Package json implements
package json

import (
	"encoding/json"
	"fmt"
)

// Unmarshal is a wrapper for json.Unmarshal
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// FormatError takes the marshalled data and the error from Unmarshal and returns a detailed
// error message where the error was and what the error is.
func FormatError(input []byte, err error) error {
	if jsonError, ok := err.(*json.SyntaxError); ok {
		line, character, offsetError := lineAndCharacter(input, int(jsonError.Offset))
		if offsetError != nil {
			return err
		}

		return fmt.Errorf("syntax error at line %d, character %d: %w", line, character, err)
	}

	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		line, character, offsetError := lineAndCharacter(input, int(jsonError.Offset))
		if offsetError != nil {
			return err
		}

		return fmt.Errorf("expect type '%s' for '%s' at line %d, character %d: %w", jsonError.Type.String(), jsonError.Field, line, character, err)
	}

	return err
}

func lineAndCharacter(input []byte, offset int) (line int, character int, err error) {
	lf := byte(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("couldn't find offset %d within the input", offset)
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	return line, character, nil
}
