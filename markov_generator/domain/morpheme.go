package domain

import "fmt"

type Surface string
type Reading string

type Morpheme struct {
	Surface Surface
	Reading Reading
}

func (m *Morpheme) String() string {
	return fmt.Sprintf("%s(%s)", m.Surface, m.Reading)
}
