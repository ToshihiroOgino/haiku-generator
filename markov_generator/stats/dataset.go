package stats

import (
	"markov_generator/domain"
	"markov_generator/mecab"

	"github.com/gookit/slog"
)

type HaikuData map[domain.Season](map[domain.Kigo]([]domain.Uta))

func (h HaikuData) ParseAll() VocabularySet {
	instance := mecab.CreateInstance()
	defer instance.Close()
	list := []*domain.Morpheme{}
	for _, season := range h {
		for _, kigo := range season {
			for _, uta := range kigo {
				if uta.IsInvalid() {
					continue
				}
				morphemes := instance.Exec(uta.GetCleaned().String())
				list = append(list, morphemes...)
			}
		}
	}
	vocab := CreateVocabularySet(list)
	slog.Debugf("parsed %d morphemes", len(list))
	slog.Debugf("got %d vocabularies", len(vocab))
	return vocab
}
