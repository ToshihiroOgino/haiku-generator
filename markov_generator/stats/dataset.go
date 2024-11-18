package stats

import (
	"markov_generator/domain"
	"markov_generator/mecab"

	"github.com/gookit/slog"
)

type HaikuData map[domain.Season](map[domain.Kigo]([]domain.Uta))

type Data struct {
	Uta       domain.Uta
	KigoPos   int
	Morphemes []domain.Morpheme
}
type Dataset map[domain.Season](map[domain.Kigo]([]domain.Uta))

func (h HaikuData) ParseAll() VocabularySet {
	instance := mecab.CreateInstance()
	defer instance.Close()
	list := []domain.Morpheme{}
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

func (h HaikuData) CreateCorpus() *Corpus {
	instance := mecab.CreateInstance()
	defer instance.Close()
	corpus := InitCorpus()
	for _, season := range h {
		slog.Infof("creating corpus from %s's haiku", season)
		for _, kigo := range season {
			for _, uta := range kigo {
				if uta.IsInvalid() {
					continue
				}
				morphemes := instance.Exec(uta.GetCleaned().String())
				corpus.Update(morphemes)
			}
		}
	}
	return corpus
}
