package domain

import (
	"fmt"
	"unicode/utf8"
)

type Surface string
type Reading string

type Morpheme struct {
	Surface Surface
	Reading Reading
}

func (m *Morpheme) String() string {
	return fmt.Sprintf("%s_%s", m.Surface, m.Reading)
}

func (s Surface) String() string {
	return string(s)
}

func (r Reading) String() string {
	return string(r)
}

func (s Surface) Length() int {
	return utf8.RuneCountInString(s.String())
}

func (r Reading) Length() int {
	return utf8.RuneCountInString(r.String())
}
