package stats

import "markov_generator/domain"

type Corpus struct {
	VocabularySet VocabularySet
	MorphemeList  []*domain.Morpheme
	bos           *Vocabulary
	eos           *Vocabulary
}

func InitCorpus() *Corpus {
	c := &Corpus{
		VocabularySet: VocabularySet{},
		MorphemeList:  []*domain.Morpheme{},
	}

	var bos = &domain.Morpheme{Surface: "BOS", Reading: "BOS"}
	var eos = &domain.Morpheme{Surface: "EOS", Reading: "EOS"}
	c.bos = c.InsertMorpheme(bos)
	c.eos = c.InsertMorpheme(eos)
	return c
}

func (c *Corpus) InsertMorpheme(m *domain.Morpheme) *Vocabulary {
	vocab, duplicate := c.VocabularySet[*m]
	if duplicate {
		return vocab
	}
	newId := VocabID(len(c.MorphemeList))
	vocab = &Vocabulary{ID: newId, Next: NextVocab{}, Prev: PrevVocab{}}
	c.VocabularySet[*m] = vocab
	c.MorphemeList = append(c.MorphemeList, m)
	return vocab
}

func (c *Corpus) GetMorphemeFromID(id VocabID) *domain.Morpheme {
	return c.MorphemeList[id]
}

func (c *Corpus) GetIDFromMorpheme(m *domain.Morpheme) VocabID {
	vocab, exists := c.VocabularySet[*m]
	if !exists {
		return -1
	}
	return vocab.ID
}

func (c *Corpus) Update(morphemes []domain.Morpheme) {
	vocabList := []*Vocabulary{c.bos}
	for _, m := range morphemes {
		vocab := c.InsertMorpheme(&m)
		vocabList = append(vocabList, vocab)
	}
	vocabList = append(vocabList, c.eos)
	for i := 1; i < len(vocabList)-1; i++ {
		prev := vocabList[i-1]
		vocabList[i].Next[prev.ID]++
		next := vocabList[i+1]
		vocabList[i].Prev[next.ID]++
	}
	c.bos.Next[vocabList[1].ID]++
	c.eos.Prev[vocabList[len(vocabList)-2].ID]++
}
