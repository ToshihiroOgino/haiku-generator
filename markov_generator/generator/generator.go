package generator

import (
	"fmt"
	"markov_generator/domain"
	"markov_generator/stats"
	"math/rand"
	"sort"
	"unicode/utf8"
)

func GenerateFromBegin(corpus *stats.Corpus) domain.Uta {
	uta := ""
	prevID := corpus.Bos.ID
	length := 0
	utaReading := ""
	for length < 5+7+5 {
		nextID := generateNext(corpus, prevID)
		if nextID == corpus.Eos.ID {
			break
		}

		m := corpus.GetMorphemeFromID(nextID)
		uta += m.Surface.String()
		length += utf8.RuneCountInString(m.Reading.String())
		utaReading += m.Reading.String()
		prevID = nextID
		// slog.Debug("uta", uta)
	}
	uta = fmt.Sprintf("%s (%s %d音)", uta, utaReading, length)
	return domain.Uta(uta)
}

func generateNext(c *stats.Corpus, prevID stats.VocabID) stats.VocabID {
	prev := c.GetVocabularyFromID(prevID)
	arr := make([]int, len(prev.Next))
	randMax := 0
	idx := 0
	for _, count := range prev.Next {
		randMax += count
		arr[idx] = randMax
		idx++
	}
	// slog.Debug("prevID", prevID)
	// slog.Debug("prev.Next", prev.Next)
	// slog.Debug("arr", arr)
	// slog.Debug("randMax", randMax)
	if randMax == 0 {
		return c.Eos.ID
	}
	rand := rand.Intn(randMax)
	idx = sort.Search(len(arr), func(i int) bool {
		return arr[i] >= rand
	})
	keys := make([]stats.VocabID, 0, len(prev.Next))
	for vi := range prev.Next {
		keys = append(keys, vi)
	}
	return stats.VocabID(keys[idx])
}
