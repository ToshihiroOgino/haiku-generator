package main

import (
	"markov_generator/fileio"
	"markov_generator/mecab"
	"strings"
	"unicode/utf8"

	"github.com/gookit/slog"
)

func fn1() {
	instance := mecab.CreateInstance()
	defer instance.Close()
	haiku := [...]string{"塵取にはこびて藍を植ゑにけり"}
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
		panic(err)
	}
	vocab := data.ParseAll()
	err = fileio.SaveVocabularySet("./fileio/vocab.json", vocab)
	if err != nil {
		panic(err)
	}
}

func corpus() {
	if data, err := fileio.LoadHaikuData("./fileio/haiku.json"); err != nil {
		panic(err)
	} else {
		corpus := data.CreateCorpus()
		if err := fileio.SaveCorpus("./fileio/corpus.json", corpus); err != nil {
			panic(err)
		}
	}
}

func kigoStat() {
	if data, err := fileio.LoadHaikuData("./fileio/haiku.json"); err != nil {
		panic(err)
	} else {
		s := data.KigoAnalyze()
		if err := fileio.SaveKigoStat("./fileio/kigo_stat.json", s); err != nil {
			panic(err)
		}
	}
}

func main() {
	kigoStat()
}
