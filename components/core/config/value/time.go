package value

import "time"

// time

type Time time.Time

func NewTime(p *time.Time, val time.Time) *Time {
	*p = val

	return (*Time)(p)
}

func (u *Time) Set(val string) error {
	v, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return err
	}
	*u = Time(v)
	return nil
}

func (u *Time) String() string {
	v := time.Time(*u)
	return v.Format(time.RFC3339)
}

func (u *Time) Validate() error {
	return nil
}

func (u *Time) IsEmpty() bool {
	v := time.Time(*u)
	return v.IsZero()
}
