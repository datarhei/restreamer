package value

import (
	"sort"
	"strconv"
	"strings"
)

// string

type String string

func NewString(p *string, val string) *String {
	*p = val

	return (*String)(p)
}

func (s *String) Set(val string) error {
	*s = String(val)
	return nil
}

func (s *String) String() string {
	return string(*s)
}

func (s *String) Validate() error {
	return nil
}

func (s *String) IsEmpty() bool {
	return len(string(*s)) == 0
}

// array of strings

type StringList struct {
	p         *[]string
	separator string
}

func NewStringList(p *[]string, val []string, separator string) *StringList {
	v := &StringList{
		p:         p,
		separator: separator,
	}

	*p = val

	return v
}

func (s *StringList) Set(val string) error {
	list := []string{}

	for _, elm := range strings.Split(val, s.separator) {
		elm = strings.TrimSpace(elm)
		if len(elm) != 0 {
			list = append(list, elm)
		}
	}

	*s.p = list

	return nil
}

func (s *StringList) String() string {
	if s.IsEmpty() {
		return "(empty)"
	}

	return strings.Join(*s.p, s.separator)
}

func (s *StringList) Validate() error {
	return nil
}

func (s *StringList) IsEmpty() bool {
	return len(*s.p) == 0
}

// map of strings to strings

type StringMapString struct {
	p *map[string]string
}

func NewStringMapString(p *map[string]string, val map[string]string) *StringMapString {
	v := &StringMapString{
		p: p,
	}

	if *p == nil {
		*p = make(map[string]string)
	}

	if val != nil {
		*p = val
	}

	return v
}

func (s *StringMapString) Set(val string) error {
	mappings := make(map[string]string)

	for _, elm := range strings.Split(val, " ") {
		elm = strings.TrimSpace(elm)
		if len(elm) == 0 {
			continue
		}

		mapping := strings.SplitN(elm, ":", 2)

		mappings[mapping[0]] = mapping[1]
	}

	*s.p = mappings

	return nil
}

func (s *StringMapString) String() string {
	if s.IsEmpty() {
		return "(empty)"
	}

	sms := *s.p

	keys := []string{}
	for k := range sms {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	mappings := make([]string, len(*s.p))

	i := 0
	for _, k := range keys {
		mappings[i] = k + ":" + sms[k]
		i++
	}

	return strings.Join(mappings, " ")
}

func (s *StringMapString) Validate() error {
	return nil
}

func (s *StringMapString) IsEmpty() bool {
	return len(*s.p) == 0
}

// boolean

type Bool bool

func NewBool(p *bool, val bool) *Bool {
	*p = val

	return (*Bool)(p)
}

func (b *Bool) Set(val string) error {
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*b = Bool(v)
	return nil
}

func (b *Bool) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *Bool) Validate() error {
	return nil
}

func (b *Bool) IsEmpty() bool {
	return !bool(*b)
}

// int

type Int int

func NewInt(p *int, val int) *Int {
	*p = val

	return (*Int)(p)
}

func (i *Int) Set(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}

func (i *Int) String() string {
	return strconv.Itoa(int(*i))
}

func (i *Int) Validate() error {
	return nil
}

func (i *Int) IsEmpty() bool {
	return int(*i) == 0
}

// int64

type Int64 int64

func NewInt64(p *int64, val int64) *Int64 {
	*p = val

	return (*Int64)(p)
}

func (u *Int64) Set(val string) error {
	v, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return err
	}
	*u = Int64(v)
	return nil
}

func (u *Int64) String() string {
	return strconv.FormatInt(int64(*u), 10)
}

func (u *Int64) Validate() error {
	return nil
}

func (u *Int64) IsEmpty() bool {
	return int64(*u) == 0
}

// uint64

type Uint64 uint64

func NewUint64(p *uint64, val uint64) *Uint64 {
	*p = val

	return (*Uint64)(p)
}

func (u *Uint64) Set(val string) error {
	v, err := strconv.ParseUint(val, 0, 64)
	if err != nil {
		return err
	}
	*u = Uint64(v)
	return nil
}

func (u *Uint64) String() string {
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *Uint64) Validate() error {
	return nil
}

func (u *Uint64) IsEmpty() bool {
	return uint64(*u) == 0
}
