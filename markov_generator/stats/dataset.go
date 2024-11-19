package stats

import (
	"markov_generator/domain"
	"markov_generator/mecab"
	"strings"

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
	for season, haikuOnKigo := range h {
		slog.Infof("creating corpus from %s's haiku", season)
		for _, kigo := range haikuOnKigo {
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

func (h HaikuData) KigoAnalyze() *KigoStat {
	instance := mecab.CreateInstance()
	defer instance.Close()

	s := &KigoStat{
		TotalAverage: 0,
		Individual:   make(map[domain.Season]([]*KigoInfo)),
	}
	numHaiku := 0
	totalSum := 0

	for season, kigoUtaMap := range h {
		slog.Infof("analyzing %s's kigo", season)
		for kigo, utaList := range kigoUtaMap {
			// 季語に対応する詩がない場合にはスキップ
			if len(utaList) == 0 {
				continue
			}
			count := 0
			sum := 0
			for _, uta := range utaList {
				if uta.IsInvalid() {
					continue
				}
				utaStr := uta.GetCleaned().String()
				idx := strings.Index(utaStr, kigo.String())
				// 季語がひらがなになっていて単純な検索が効かない場合にはスキップ (e.g. 時雨->しぐる)
				if idx == -1 {
					continue
				}
				// 季語までをMecabで読みに変換
				morphemes := instance.Exec(utaStr[:idx])
				pos := 0
				// 形態素リストの音素数を数える
				for _, m := range morphemes {
					pos += m.Surface.Length()
				}
				count++
				sum += pos
			}
			// 有効な詩がない場合にもスキップ
			if count == 0 {
				continue
			}
			s.Individual[season] = append(s.Individual[season], &KigoInfo{
				Kigo:       kigo,
				NumHaiku:   count,
				AveragePos: KigoAveragePos(sum) / KigoAveragePos(count),
			})
			numHaiku += count
			totalSum += sum
		}
	}
	s.TotalAverage = KigoAveragePos(totalSum) / KigoAveragePos(numHaiku)
	return s
}
