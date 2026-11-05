package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPattern(t *testing.T) {
	p := NewPattern("bla", "label1", "value1", "label2")
	require.Equal(t, false, p.IsValid())

	p = NewPattern("bla", "label1", "value1", "label2", "valu(e2")
	require.Equal(t, false, p.IsValid())

	p = NewPattern("bla")
	require.Equal(t, true, p.IsValid())
	require.Equal(t, "bla", p.Name())

	p = NewPattern("bla", "label1", "value1", "label2", "value2")
	require.Equal(t, true, p.IsValid())
}

func TestPatternMatch(t *testing.T) {
	p := NewPattern("bla", "label1", "value1", "label2")
	require.Equal(t, false, p.IsValid())
	require.Equal(t, false, p.Match(map[string]string{"label1": "value1"}))

	p0 := NewPattern("bla")
	require.Equal(t, true, p0.IsValid())
	require.Equal(t, true, p0.Match(map[string]string{}))
	require.Equal(t, true, p0.Match(map[string]string{"labelX": "foobar"}))

	p = NewPattern("bla", "label1", "value.", "label2", "val?ue2")
	require.Equal(t, true, p.IsValid())
	require.Equal(t, false, p.Match(map[string]string{}))
	require.Equal(t, false, p.Match(map[string]string{"label1": "value1"}))
	require.Equal(t, true, p.Match(map[string]string{"label1": "value1", "label2": "value2"}))
	require.Equal(t, true, p.Match(map[string]string{"label1": "value5", "label2": "vaue2"}))
}

func TestValue(t *testing.T) {
	d := NewDesc("group", "", []string{"label1", "label2"})
	v := NewValue(d, 42, "foobar")
	require.Nil(t, v)

	v = NewValue(d, 42, "foobar", "foobaz")
	require.NotNil(t, v)
	require.Equal(t, float64(42), v.Val())

	require.Equal(t, "", v.L("labelX"))
	require.Equal(t, "foobar", v.L("label1"))
	require.Equal(t, "foobaz", v.L("label2"))
	require.Equal(t, "group", v.Name())
	require.Equal(t, "group:label1=foobar label2=foobaz ", v.Hash())
	require.Equal(t, "group: 42.000000 {label1=foobar label2=foobaz }", v.String())

	require.Equal(t, map[string]string{"label1": "foobar", "label2": "foobaz"}, v.Labels())
}

func TestValuePattern(t *testing.T) {
	d := NewDesc("group", "", []string{"label1", "label2"})
	v := NewValue(d, 42, "foobar", "foobaz")

	p1 := NewPattern("group")
	p2 := NewPattern("group", "label1", "foobar")
	p3 := NewPattern("group", "label2", "foobaz")
	p4 := NewPattern("group", "label2", "foobaz", "label1", "foobar")

	require.Equal(t, true, v.Match(nil))
	require.Equal(t, true, v.Match([]Pattern{p1}))
	require.Equal(t, true, v.Match([]Pattern{p2}))
	require.Equal(t, true, v.Match([]Pattern{p3}))
	require.Equal(t, true, v.Match([]Pattern{p4}))
	require.Equal(t, true, v.Match([]Pattern{p1, p2, p3, p4}))

	p5 := NewPattern("group", "label1", "foobaz")

	require.Equal(t, false, v.Match([]Pattern{p5}))

	require.Equal(t, true, v.Match([]Pattern{p4, p5}))
	require.Equal(t, true, v.Match([]Pattern{p5, p4}))
}

func TestDescription(t *testing.T) {
	d := NewDesc("name", "blabla", []string{"label"})

	require.Equal(t, "name", d.Name())
	require.Equal(t, "blabla", d.Description())
	require.ElementsMatch(t, []string{"label"}, d.Labels())
	require.Equal(t, "name: blabla (label)", d.String())
}

func TestMetrics(t *testing.T) {
	m := NewMetrics()

	require.Equal(t, "", m.String())
	require.Equal(t, 0, len(m.All()))

	d := NewDesc("group", "", []string{"label1", "label2"})
	v1 := NewValue(d, 42, "foobar", "foobaz")
	require.NotNil(t, v1)

	m.Add(v1)

	require.Equal(t, v1.String(), m.String())
	require.Equal(t, 1, len(m.All()))

	l := m.Labels("group", "label2")

	require.ElementsMatch(t, []string{"foobaz"}, l)

	v2 := NewValue(d, 77, "barfoo", "bazfoo")

	m.Add(v2)

	require.Equal(t, v1.String()+v2.String(), m.String())
	require.Equal(t, 2, len(m.All()))

	l = m.Labels("group", "label2")

	require.ElementsMatch(t, []string{"foobaz", "bazfoo"}, l)

	v := m.Value("bla", "label1", "foo*")

	require.Equal(t, nullValue, v)

	v = m.Value("group")

	require.NotEqual(t, nullValue, v)

	v = m.Value("group", "label1", "foo*")

	require.NotEqual(t, nullValue, v)

	v = m.Value("group", "label2", "baz")

	require.NotEqual(t, nullValue, v)

	vs := m.Values("group")

	require.Equal(t, 2, len(vs))

	vs = m.Values("group", "label1", "foo*")

	require.Equal(t, 2, len(vs))

	vs = m.Values("group", "label2", "*baz*")

	require.NotEqual(t, 2, len(vs))

	vs = m.Values("group", "label1")

	require.Equal(t, 0, len(vs))
}
