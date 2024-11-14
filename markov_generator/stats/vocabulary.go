package stats

import (
	"fmt"
	. "markov_generator/domain"
)

type VocabID int

type Vocabulary struct {
	ID VocabID
	Reading
}

func CreateVocabularySet(data []*Morpheme) (map[Surface]Vocabulary, error) {
	id := VocabID(0)
	vocab := map[Surface]Vocabulary{}
	for _, m := range data {
		val, duplicate := vocab[m.Surface]
		if !duplicate {
			vocab[m.Surface] = Vocabulary{ID: id, Reading: m.Reading}
			id++
		} else if val.Reading != m.Reading {
			return nil, fmt.Errorf("duplicate surface with different reading: %s -> (%s, %s)", m.Surface, val.Reading, m.Reading)
		}
	}
	return vocab, nil
}
