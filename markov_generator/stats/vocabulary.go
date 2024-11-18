package stats

import (
	. "markov_generator/domain"
)

type VocabID int

type Vocabulary struct {
	ID VocabID
}

type VocabularySet map[Morpheme]Vocabulary

func CreateVocabularySet(data []*Morpheme) VocabularySet {
	id := VocabID(0)
	vocab := VocabularySet{}
	for _, m := range data {
		_, duplicate := vocab[*m]
		if !duplicate {
			vocab[*m] = Vocabulary{ID: id}
			id++
		} else {
			continue
		}
	}
	return vocab
}
