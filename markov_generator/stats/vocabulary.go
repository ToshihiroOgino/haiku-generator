package stats

import "markov_generator/domain"

type VocabID int
type Count int
type NextVocab map[VocabID]Count
type PrevVocab map[VocabID]Count
type Vocabulary struct {
	ID   VocabID
	Next NextVocab
	Prev PrevVocab
}

type VocabularySet map[domain.Morpheme]*Vocabulary

func CreateVocabularySet(data []domain.Morpheme) VocabularySet {
	id := VocabID(0)
	vocab := VocabularySet{}
	for _, m := range data {
		_, duplicate := vocab[m]
		if !duplicate {
			vocab[m] = &Vocabulary{ID: id}
			id++
		} else {
			continue
		}
	}
	return vocab
}
