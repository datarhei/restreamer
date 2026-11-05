package scalars

import (
	"fmt"
	"io"
)

type Uint64 uint64

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (u *Uint64) UnmarshalGQL(v interface{}) error {
	value, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("Uint64 must be a uint64")
	}

	*u = Uint64(value)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (u Uint64) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, "%d", uint64(u))
}
