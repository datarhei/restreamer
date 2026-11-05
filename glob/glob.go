package glob

import (
	"github.com/gobwas/glob"
)

type Glob interface {
	Match(name string) bool
}

type globber struct {
	glob glob.Glob
}

func Compile(pattern string, separators ...rune) (Glob, error) {
	g, err := glob.Compile(pattern, separators...)
	if err != nil {
		return nil, err
	}

	return &globber{glob: g}, nil
}

func (g *globber) Match(name string) bool {
	return g.glob.Match(name)
}

// Match returns whether the name matches the glob pattern, also considering
// one or several optionnal separator. An error is only returned if the pattern
// is invalid.
func Match(pattern, name string, separators ...rune) (bool, error) {
	g, err := glob.Compile(pattern, separators...)
	if err != nil {
		return false, err
	}

	return g.Match(name), nil
}
