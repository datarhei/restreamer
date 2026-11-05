package ffmpeg

import (
	"fmt"
	"regexp"
	"strings"
)

// The Validator interface is for validating strings whether they are eligible
// as input or output for a ffmpeg process
type Validator interface {
	// IsValid tests whether a text matches any expression of inputs
	IsValid(text string) bool
}

type validator struct {
	allow []*regexp.Regexp // List of compiled input expressions
	block []*regexp.Regexp // List of compiled output expressions
}

// NewValidator creates a new Validator with the given input and output
// expressions. Empty expressions will be ignored. Returns an
// error if an expression can't be compiled.
func NewValidator(allow, block []string) (Validator, error) {
	v := &validator{}

	for _, expression := range allow {
		expression = strings.TrimSpace(expression)
		if len(expression) == 0 {
			continue
		}

		re, err := regexp.Compile(expression)
		if err != nil {
			return nil, fmt.Errorf("invalid allow expression: '%s' (%w)", expression, err)
		}

		v.allow = append(v.allow, re)
	}

	for _, expression := range block {
		expression = strings.TrimSpace(expression)
		if len(expression) == 0 {
			continue
		}

		re, err := regexp.Compile(expression)
		if err != nil {
			return nil, fmt.Errorf("invalid block expression: '%s' (%w)", expression, err)
		}

		v.block = append(v.block, re)
	}

	return v, nil
}

func (v *validator) IsValid(text string) bool {
	for _, e := range v.block {
		if e.MatchString(text) {
			return false
		}
	}

	if len(v.allow) == 0 {
		return true
	}

	for _, e := range v.allow {
		if e.MatchString(text) {
			return true
		}
	}

	return false
}
