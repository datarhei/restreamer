package scalars

import (
	"fmt"
	"io"
	"time"
)

type MetricsResponseValue struct {
	TS     time.Time
	Value  float64
	isNull bool
}

func (u *MetricsResponseValue) UnmarshalGQL(v interface{}) error {
	*u = MetricsResponseValue{
		TS:     time.Time{},
		Value:  0,
		isNull: false,
	}

	return nil
}

// MarshalJSON marshals a MetricsResponseValue to GQL
func (v MetricsResponseValue) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, "[%d,", v.TS.Unix())

	if v.isNull {
		fmt.Fprintf(w, "null")
	} else {
		if v.Value == float64(int64(v.Value)) {
			fmt.Fprintf(w, "%.0f", v.Value) // 0 decimal if integer
		} else {
			fmt.Fprintf(w, "%.3f", v.Value) // max. 3 decimal if float
		}
	}

	fmt.Fprintf(w, "]")
}
