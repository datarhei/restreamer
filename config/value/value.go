package value

type Value interface {
	// String returns a string representation of the value.
	String() string

	// Set a new value for the value. Returns an
	// error if the given string representation can't
	// be transformed to the value. Returns nil
	// if the new value has been set.
	Set(string) error

	// Validate the value. The returned error will
	// indicate what is wrong with the current value.
	// Returns nil if the value is OK.
	Validate() error

	// IsEmpty returns whether the value represents an empty
	// representation for that value.
	IsEmpty() bool
}
