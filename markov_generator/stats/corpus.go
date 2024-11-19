package stats

import (
	"encoding/json"
	"markov_generator/domain"
)

type Corpus struct {
	VocabularySet VocabularySet
	MorphemeList  []*domain.Morpheme
	Bos           *Vocabulary
	Eos           *Vocabulary
}

func InitCorpus() *Corpus {
	c := &Corpus{
		VocabularySet: VocabularySet{},
		MorphemeList:  []*domain.Morpheme{},
	}

	var bos = &domain.Morpheme{Surface: "BOS", Reading: "BOS"}
	var eos = &domain.Morpheme{Surface: "EOS", Reading: "EOS"}
	c.Bos = c.InsertMorpheme(bos)
	c.Eos = c.InsertMorpheme(eos)
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

func (c *Corpus) GetVocabularyFromID(id VocabID) *Vocabulary {
	return c.VocabularySet[*c.GetMorphemeFromID(id)]
}

func (c *Corpus) Update(morphemes []domain.Morpheme) {
	vocabList := []*Vocabulary{c.Bos}
	for _, m := range morphemes {
		vocab := c.InsertMorpheme(&m)
		vocabList = append(vocabList, vocab)
	}
	vocabList = append(vocabList, c.Eos)
	for i := 1; i < len(vocabList)-1; i++ {
		if i != 1 {
			prev := vocabList[i-1]
			vocabList[i].Prev[prev.ID]++
		}
		if i != len(vocabList)-2 {
			next := vocabList[i+1]
			vocabList[i].Next[next.ID]++
		}
	}
	c.Bos.Next[vocabList[1].ID]++
	c.Eos.Prev[vocabList[len(vocabList)-2].ID]++
}

/* ----- For JSON marshalling and unmarshalling ----- */
type VocabMarshalItem struct {
	Surface domain.Surface
	Reading domain.Reading
	Next    map[VocabID]int
	Prev    map[VocabID]int
}

func (c Corpus) MarshalJSON() ([]byte, error) {
	vocab := make(map[VocabID]VocabMarshalItem)
	for k, v := range c.VocabularySet {
		vocab[v.ID] = VocabMarshalItem{
			Surface: k.Surface,
			Reading: k.Reading,
			Next:    v.Next,
			Prev:    v.Prev,
		}
	}
	return json.Marshal(struct {
		Vocabulary map[VocabID]VocabMarshalItem
		BosID      VocabID
		EosID      VocabID
	}{
		Vocabulary: vocab,
		BosID:      c.Bos.ID,
		EosID:      c.Eos.ID,
	})
}

func (c *Corpus) UnmarshalJSON(data []byte) error {
	var obj struct {
		Vocabulary map[VocabID]VocabMarshalItem
		BosID      VocabID
		EosID      VocabID
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	c.VocabularySet = make(VocabularySet)
	c.MorphemeList = make([]*domain.Morpheme, len(obj.Vocabulary))
	for id, item := range obj.Vocabulary {
		morpheme := &domain.Morpheme{
			Surface: item.Surface,
			Reading: item.Reading,
		}
		vocab := &Vocabulary{
			ID:   id,
			Next: item.Next,
			Prev: item.Prev,
		}
		c.VocabularySet[*morpheme] = vocab
		c.MorphemeList[id] = morpheme
		if item.Surface == "BOS" {
			c.Bos = vocab
		} else if item.Surface == "EOS" {
			c.Eos = vocab
		}
	}
	return nil
}
