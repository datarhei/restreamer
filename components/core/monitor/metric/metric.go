package metric

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Pattern interface {
	// Name returns the name of the metric this pattern is designated to.
	Name() string

	// Match returns whether a map of labels with its label values
	// match this pattern. All labels have to be present and need to match.
	Match(labels map[string]string) bool

	// IsValid returns whether the pattern is valid.
	IsValid() bool
}

type pattern struct {
	name   string
	labels map[string]*regexp.Regexp
	valid  bool
}

// NewPattern creates a new pattern with the given prefix and group name. There
// has to be an even number of labels, which is ("label", "labelvalue", "label",
// "labelvalue" ...). The label value will be interpreted as regular expression.
func NewPattern(name string, labels ...string) Pattern {
	p := &pattern{
		name:   name,
		labels: make(map[string]*regexp.Regexp),
	}

	if len(labels)%2 == 0 {
		for i := 0; i < len(labels); i += 2 {
			exp, err := regexp.Compile(labels[i+1])
			if err != nil {
				continue
			}

			p.labels[labels[i]] = exp
		}
	}

	if len(p.labels) == len(labels)/2 {
		p.valid = true
	}

	return p
}

func (p *pattern) Name() string {
	return p.name
}

func (p *pattern) Match(labels map[string]string) bool {
	if !p.valid {
		return false
	}

	if len(p.labels) == 0 {
		return true
	}

	for pname, pexp := range p.labels {
		value, ok := labels[pname]
		if !ok {
			return false
		}

		if !pexp.MatchString(value) {
			return false
		}
	}

	return true
}

func (p *pattern) IsValid() bool {
	return p.valid
}

// Metrics is a collection of values
type Metrics interface {
	// Value returns the first value that matches the name and the labels. The labels
	// are used to create a pattern and therefore must obey to the rules of NewPattern.
	Value(name string, labels ...string) Value

	// Values returns all values that matches the name and the labels. The labels
	// are used to create a pattern and therefore must obey to the rules of NewPattern.
	Values(name string, labels ...string) []Value

	// Labels return a list of all values for a label.
	Labels(name string, label string) []string

	// All returns all values currently stored in the collection.
	All() []Value

	// Add adds a value to the collection.
	Add(v Value)

	// String return a string representation of all collected values.
	String() string
}

// metrics is an implementation of the Metrics interface.
type metrics struct {
	values []Value
}

// NewMetrics returns a new metrics instance.
func NewMetrics() *metrics {
	return &metrics{}
}

func (m *metrics) String() string {
	s := ""

	for _, v := range m.values {
		s += v.String()
	}

	return s
}

func (m *metrics) Values(name string, labels ...string) []Value {
	if len(labels)%2 != 0 {
		return []Value{}
	}

	patterns := []Pattern{
		NewPattern(name, labels...),
	}

	values := []Value{}

	for _, v := range m.values {
		if !v.Match(patterns) {
			continue
		}

		values = append(values, v)
	}

	return values
}

func (m *metrics) Value(name string, labels ...string) Value {
	vs := m.Values(name, labels...)

	if len(vs) == 0 {
		return nullValue
	}

	return vs[0]
}

func (m *metrics) All() []Value {
	return m.values
}

func (m *metrics) Labels(name string, label string) []string {
	vs := m.Values(name)
	values := make(map[string]struct{})

	for _, v := range vs {
		l := v.L(label)
		if len(l) == 0 {
			break
		}

		values[l] = struct{}{}
	}

	labelvalues := []string{}

	for v := range values {
		labelvalues = append(labelvalues, v)
	}

	return labelvalues
}

func (m *metrics) Add(v Value) {
	if v == nil {
		return
	}

	m.values = append(m.values, v)
}

type Value interface {
	Name() string
	Val() float64
	L(name string) string
	Labels() map[string]string
	Match(patterns []Pattern) bool
	Hash() string
	String() string
}

type value struct {
	name   string
	labels map[string]string
	value  float64
	hash   string
}

var nullValue Value = NewValue(NewDesc("", "", nil), 0)

func NewValue(description *Description, v float64, elms ...string) Value {
	if len(description.labels) != len(elms) {
		return nil
	}

	val := &value{
		name:   description.name,
		value:  v,
		labels: make(map[string]string),
	}

	labels := []string{}

	for i, label := range description.labels {
		val.labels[label] = elms[i]
		labels = append(labels, label)
	}

	val.hash = fmt.Sprintf("%s:", val.name)

	sort.Strings(labels)
	for _, k := range labels {
		val.hash += k + "=" + val.labels[k] + " "
	}

	return val
}

func (v *value) Hash() string {
	return v.hash
}

func (v *value) String() string {
	s := fmt.Sprintf("%s: %f {", v.name, v.value)

	keys := []string{}
	for k := range v.labels {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += k + "=" + v.labels[k] + " "
	}

	s += "}"

	return s
}

func (v *value) Name() string {
	return v.name
}

func (v *value) Val() float64 {
	return v.value
}

func (v *value) L(name string) string {
	l, ok := v.labels[name]
	if !ok {
		return ""
	}

	return l
}

func (v *value) Labels() map[string]string {
	if len(v.labels) == 0 {
		return nil
	}

	l := make(map[string]string, len(v.labels))

	for k, v := range v.labels {
		l[k] = v
	}

	return l
}

func (v *value) Match(patterns []Pattern) bool {
	if len(patterns) == 0 {
		return true
	}

	for _, p := range patterns {
		if v.name != p.Name() {
			continue
		}

		if !p.Match(v.labels) {
			continue
		}

		return true
	}

	return false
}

type Description struct {
	name        string
	description string
	labels      []string
}

func NewDesc(name, description string, labels []string) *Description {
	return &Description{
		name:        name,
		description: description,
		labels:      labels,
	}
}

func (d *Description) String() string {
	return fmt.Sprintf("%s: %s (%s)", d.name, d.description, strings.Join(d.labels, ","))
}

func (d *Description) Name() string {
	return d.name
}

func (d *Description) Description() string {
	return d.description
}

func (d *Description) Labels() []string {
	labels := make([]string, len(d.labels))
	copy(labels, d.labels)

	return labels
}

type Collector interface {
	Prefix() string
	Describe() []*Description
	Collect() Metrics
	Stop()
}

type Reader interface {
	Collect([]Pattern) Metrics
}
