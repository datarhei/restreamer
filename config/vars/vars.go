package vars

import (
	"fmt"
	"os"

	"github.com/datarhei/core/v16/config/value"
)

type variable struct {
	value       value.Value // The actual value
	defVal      string      // The default value in string representation
	name        string      // A name for this value
	envName     string      // The environment variable that corresponds to this value
	envAltNames []string    // Alternative environment variable names
	description string      // A desriptions for this value
	required    bool        // Whether a non-empty value is required
	disguise    bool        // Whether the value should be disguised if printed
	merged      bool        // Whether this value has been replaced by its corresponding environment variable
}

type Variable struct {
	Value       string
	Name        string
	EnvName     string
	Description string
	Merged      bool
}

type message struct {
	message  string   // The log message
	variable Variable // The config field this message refers to
	level    string   // The loglevel for this message
}

type Variables struct {
	vars []*variable
	logs []message
}

func (vs *Variables) Register(val value.Value, name, envName string, envAltNames []string, description string, required, disguise bool) {
	vs.vars = append(vs.vars, &variable{
		value:       val,
		defVal:      val.String(),
		name:        name,
		envName:     envName,
		envAltNames: envAltNames,
		description: description,
		required:    required,
		disguise:    disguise,
	})
}

func (vs *Variables) Transfer(vss *Variables) {
	for _, v := range vs.vars {
		if vss.IsMerged(v.name) {
			v.merged = true
		}
	}
}

func (vs *Variables) SetDefault(name string) {
	v := vs.findVariable(name)
	if v == nil {
		return
	}

	v.value.Set(v.defVal)
}

func (vs *Variables) Get(name string) (string, error) {
	v := vs.findVariable(name)
	if v == nil {
		return "", fmt.Errorf("variable not found")
	}

	return v.value.String(), nil
}

func (vs *Variables) Set(name, val string) error {
	v := vs.findVariable(name)
	if v == nil {
		return fmt.Errorf("variable not found")
	}

	return v.value.Set(val)
}

func (vs *Variables) Log(level, name string, format string, args ...interface{}) {
	v := vs.findVariable(name)
	if v == nil {
		return
	}

	variable := Variable{
		Value:       v.value.String(),
		Name:        v.name,
		EnvName:     v.envName,
		Description: v.description,
		Merged:      v.merged,
	}

	if v.disguise {
		variable.Value = "***"
	}

	l := message{
		message:  fmt.Sprintf(format, args...),
		variable: variable,
		level:    level,
	}

	vs.logs = append(vs.logs, l)
}

func (vs *Variables) Merge() {
	for _, v := range vs.vars {
		if len(v.envName) == 0 {
			continue
		}

		var envval string
		var ok bool

		envval, ok = os.LookupEnv(v.envName)
		if !ok {
			foundAltName := false

			for _, envName := range v.envAltNames {
				envval, ok = os.LookupEnv(envName)
				if ok {
					foundAltName = true
					vs.Log("warn", v.name, "deprecated name, please use %s", v.envName)
					break
				}
			}

			if !foundAltName {
				continue
			}
		}

		err := v.value.Set(envval)
		if err != nil {
			vs.Log("error", v.name, "%s", err.Error())
		}

		v.merged = true
	}
}

func (vs *Variables) IsMerged(name string) bool {
	v := vs.findVariable(name)
	if v == nil {
		return false
	}

	return v.merged
}

func (vs *Variables) Validate() {
	for _, v := range vs.vars {
		vs.Log("info", v.name, "%s", "")

		err := v.value.Validate()
		if err != nil {
			vs.Log("error", v.name, "%s", err.Error())
		}

		if v.required && v.value.IsEmpty() {
			vs.Log("error", v.name, "a value is required")
		}
	}
}

func (vs *Variables) ResetLogs() {
	vs.logs = nil
}

func (vs *Variables) Messages(logger func(level string, v Variable, message string)) {
	for _, l := range vs.logs {
		logger(l.level, l.variable, l.message)
	}
}

func (vs *Variables) HasErrors() bool {
	for _, l := range vs.logs {
		if l.level == "error" {
			return true
		}
	}

	return false
}

func (vs *Variables) Overrides() []string {
	overrides := []string{}

	for _, v := range vs.vars {
		if v.merged {
			overrides = append(overrides, v.name)
		}
	}

	return overrides
}

func (vs *Variables) findVariable(name string) *variable {
	for _, v := range vs.vars {
		if v.name == name {
			return v
		}
	}

	return nil
}
