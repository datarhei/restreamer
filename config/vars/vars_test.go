package vars

import (
	"os"
	"testing"

	"github.com/datarhei/core/v16/config/value"

	"github.com/stretchr/testify/require"
)

func TestVars(t *testing.T) {
	v1 := Variables{}

	s := ""

	v1.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, false)

	require.Equal(t, "foobar", s)
	x, _ := v1.Get("string")
	require.Equal(t, "foobar", x)

	v := v1.findVariable("string")
	v.value.Set("barfoo")

	require.Equal(t, "barfoo", s)
	x, _ = v1.Get("string")
	require.Equal(t, "barfoo", x)

	v1.Set("string", "foobaz")

	require.Equal(t, "foobaz", s)
	x, _ = v1.Get("string")
	require.Equal(t, "foobaz", x)

	v1.SetDefault("string")

	require.Equal(t, "foobar", s)
	x, _ = v1.Get("string")
	require.Equal(t, "foobar", x)
}

func TestSetDefault(t *testing.T) {
	v := Variables{}
	s := ""

	v.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, false)

	require.Equal(t, "foobar", s)

	v.Set("string", "foobaz")

	require.Equal(t, "foobaz", s)

	v.SetDefault("strong")

	require.Equal(t, "foobaz", s)

	v.SetDefault("string")

	require.Equal(t, "foobar", s)
}

func TestGet(t *testing.T) {
	v := Variables{}

	s := ""

	v.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, false)

	value, err := v.Get("string")
	require.NoError(t, err)
	require.Equal(t, "foobar", value)

	value, err = v.Get("strong")
	require.Error(t, err)
	require.Equal(t, "", value)
}

func TestSet(t *testing.T) {
	v := Variables{}

	s := ""

	v.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, false)

	err := v.Set("string", "foobaz")
	require.NoError(t, err)
	require.Equal(t, "foobaz", s)

	err = v.Set("strong", "fooboz")
	require.Error(t, err)
	require.Equal(t, "foobaz", s)
}

func TestLog(t *testing.T) {
	v := Variables{}

	s := ""

	v.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, false)

	v.Log("info", "string", "hello %s", "world")
	require.Equal(t, 1, len(v.logs))

	v.Log("info", "strong", "hello %s", "world")
	require.Equal(t, 1, len(v.logs))

	require.Equal(t, "hello world", v.logs[0].message)
	require.Equal(t, "info", v.logs[0].level)
	require.Equal(t, Variable{
		Value:       "foobar",
		Name:        "string",
		EnvName:     "",
		Description: "a string",
		Merged:      false,
	}, v.logs[0].variable)

	v.ResetLogs()

	require.Equal(t, 0, len(v.logs))
}

func TestMerge(t *testing.T) {
	v := Variables{}

	s := ""
	os.Setenv("CORE_TEST_STRING", "foobaz")

	v.Register(value.NewString(&s, "foobar"), "string", "CORE_TEST_STRING", nil, "a string", false, false)

	require.Equal(t, s, "foobar")

	v.Merge()

	require.Equal(t, s, "foobaz")
	require.Equal(t, true, v.IsMerged("string"))
	require.Equal(t, 0, len(v.logs))

	os.Unsetenv("CORE_TEST_STRING")
}

func TestMergeAlt(t *testing.T) {
	v := Variables{}

	s := ""
	os.Setenv("CORE_TEST_STRING", "foobaz")

	v.Register(value.NewString(&s, "foobar"), "string", "CORE_TEST_STRUNG", []string{"CORE_TEST_STRING"}, "a string", false, false)

	require.Equal(t, s, "foobar")

	v.Merge()

	require.Equal(t, s, "foobaz")
	require.Equal(t, true, v.IsMerged("string"))
	require.Equal(t, 1, len(v.logs))

	require.Contains(t, v.logs[0].message, "CORE_TEST_STRUNG")
	require.Equal(t, "warn", v.logs[0].level)

	os.Unsetenv("CORE_TEST_STRING")
}

func TestNoMerge(t *testing.T) {
	v := Variables{}

	s := ""
	os.Setenv("CORE_TEST_STRONG", "foobaz")

	v.Register(value.NewString(&s, "foobar"), "string", "CORE_TEST_STRING", nil, "a string", false, false)

	require.Equal(t, s, "foobar")

	v.Merge()

	require.Equal(t, s, "foobar")
	require.Equal(t, false, v.IsMerged("string"))

	os.Unsetenv("CORE_TEST_STRONG")
}

func TestValidate(t *testing.T) {
	v := Variables{}

	s1 := ""
	s2 := ""

	v.Register(value.NewString(&s1, ""), "string", "", nil, "a string", false, false)
	v.Register(value.NewString(&s2, ""), "string", "", nil, "a string", true, false)

	require.Equal(t, s1, "")
	require.Equal(t, s2, "")

	require.Equal(t, false, v.HasErrors())

	v.Validate()

	require.Equal(t, true, v.HasErrors())

	ninfo := 0
	nerror := 0
	v.Messages(func(level string, v Variable, message string) {
		if level == "info" {
			ninfo++
		} else if level == "error" {
			nerror++
		}
	})

	require.Equal(t, 2, ninfo)
	require.Equal(t, 1, nerror)
}

func TestOverrides(t *testing.T) {
	v := Variables{}

	s := ""
	os.Setenv("CORE_TEST_STRING", "foobaz")

	v.Register(value.NewString(&s, "foobar"), "string", "CORE_TEST_STRING", nil, "a string", false, false)
	v.Merge()

	overrides := v.Overrides()

	require.ElementsMatch(t, []string{"string"}, overrides)
}

func TestDisquise(t *testing.T) {
	v := Variables{}

	s := ""

	v.Register(value.NewString(&s, "foobar"), "string", "", nil, "a string", false, true)

	v.Log("info", "string", "hello %s", "world")
	require.Equal(t, 1, len(v.logs))

	require.Equal(t, "hello world", v.logs[0].message)
	require.Equal(t, "info", v.logs[0].level)
	require.Equal(t, Variable{
		Value:       "***",
		Name:        "string",
		EnvName:     "",
		Description: "a string",
		Merged:      false,
	}, v.logs[0].variable)
}
