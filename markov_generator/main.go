package main

import (
	"encoding/json"
	"markov_generator/domain"
	"markov_generator/fileio"
	"markov_generator/generator"
	"markov_generator/mecab"
	"markov_generator/stats"
	"math/rand"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gookit/slog"
)

func fn1() {
	instance := mecab.CreateInstance()
	defer instance.Close()
	haiku := [...]string{"明日葉や終りのなきものはなし", "明日葉"}
	for _, h := range haiku {
		res := instance.Exec(h)
		// res := instance.Parse(h)
		slog.Debug("input", h)
		slog.Debug("res", res)
	}
}

func fn2() {
	str := "あの頃は輝いていた賀状書く"
	kigo := "賀状"
	idx := strings.Index(str, kigo)
	slog.Debug("idx", idx)
	instance := mecab.CreateInstance()
	defer instance.Close()
	str = str[:idx]
	slog.Debug("str", str)
	str = "あああ"
	slog.Debug("len(\"あああ\")", utf8.RuneCountInString(str))
	// arr := instance.Exec(strings.Trim(str, kigo))
}

func vocabList() {
	data, err := fileio.LoadHaikuData("./fileio/haiku.json")
	if err != nil {
		slog.FatalErr(err)
	}
	vocab := data.ParseAll()
	err = fileio.SaveVocabularySet("./fileio/vocab.json", vocab)
	if err != nil {
		slog.FatalErr(err)
	}
}

func corpus() {
	if data, err := fileio.LoadHaikuData("./fileio/haiku.json"); err != nil {
		slog.FatalErr(err)
	} else {
		corpus := data.CreateCorpus()
		if err := fileio.SaveCorpus("./fileio/corpus.json", corpus); err != nil {
			slog.FatalErr(err)
		}
	}
}

func kigoStat() {
	if data, err := fileio.LoadHaikuData("./fileio/haiku.json"); err != nil {
		slog.FatalErr(err)
	} else {
		s := data.KigoAnalyze()
		if err := fileio.SaveKigoStat("./fileio/kigo_stat.json", s); err != nil {
			slog.FatalErr(err)
		}
	}
}

func utaFromBegin() {
	if corpus, err := fileio.LoadCorpus("./fileio/corpus.json"); err != nil {
		slog.FatalErr(err)
	} else {
		for i := 0; i < 5; i++ {
			uta := generator.GenerateFromBegin(corpus)
			slog.Debug("uta", uta)
		}
	}
}

func utaFromKigo() {
	corpus, err := fileio.LoadCorpus("./fileio/corpus.json")
	if err != nil {
		slog.FatalErr(err)
	}
	kigoStat, err := fileio.LoadKigoStat("./fileio/kigo_stat.json")
	if err != nil {
		slog.FatalErr(err)
	}

	kigoList := []domain.Kigo{"蟻", "啓蟄", "花野", "苗代", "藤", "大寒", "花", "八月", "八月", "炎天"}

	slog.Debug("kigoList", kigoList, "len(kigoList)", len(kigoList))
	rand.NewSource(0)
	generated := make(stats.HaikuData)
	season := domain.Season("Unknown")
	generated[season] = make(map[domain.Kigo]([]domain.Uta))

	instance := mecab.CreateInstance()
	defer instance.Close()

	for _, kigo := range kigoList {
		slog.Infof("----- 季語: %s -----", kigo)
		generated[season][kigo] = []domain.Uta{}
		for i := 0; i < 5; i++ {
			uta := generator.GenerateFromKigo(corpus, kigoStat, kigo, instance)
			slog.Debug("uta", uta)
			generated[season][kigo] = append(generated[season][kigo], uta.Uta)
		}
	}
}

func small() {
	if data, err := fileio.LoadHaikuData("./fileio/small.json"); err != nil {
		slog.FatalErr(err)
	} else {
		corpus := data.CreateCorpus()
		// slog.Debug("corpus", corpus)
		fileio.SaveCorpus("./fileio/small_corpus.json", corpus)
		uta := generator.GenerateFromBegin(corpus)
		slog.Debug("uta", uta)
	}
}

func kigoStatAnalysis() {
	kigoStat, err := fileio.LoadKigoStat("./fileio/kigo_stat.json")
	if err != nil {
		slog.FatalErr(err)
	}
	allInfo := []*stats.KigoInfo{}
	for _, infoList := range kigoStat.Individual {
		allInfo = append(allInfo, infoList...)
	}
	slog.Debug("len(allInfo)", len(allInfo))
	sort.Slice(allInfo, func(i, j int) bool {
		return allInfo[i].NumHaiku >= allInfo[j].NumHaiku
	})
	bytes, _ := json.MarshalIndent(allInfo[:15], "", "  ")
	slog.Debug("top", string(bytes))
}

func generateALot() {
	corpus, err := fileio.LoadCorpus("./fileio/corpus.json")
	if err != nil {
		slog.FatalErr(err)
	}
	kigoStat, err := fileio.LoadKigoStat("./fileio/kigo_stat.json")
	if err != nil {
		slog.FatalErr(err)
	}

	data := generator.GenerateALot(corpus, kigoStat)

	if err := fileio.SaveGeneratedHaiku("./fileio/generated.json", data); err != nil {
		slog.FatalErr(err)
	}
}

func main() {
	// kigoStatAnalysis()
	generateALot()
	// corpus()
}
