package main

import (
	"markov_generator/fileio"
	"markov_generator/generator"
	"markov_generator/mecab"
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
	for i := 0; i < 5; i++ {
		uta := generator.GenerateFromKigo(corpus, kigoStat, "冬田")
		slog.Debug("uta", uta)
	}
}

func bin() {
	src := make([]int, 10)
	for i := 0; i < 10; i++ {
		src[i] = rand.Intn(30)
	}
	slog.Debug("src", src)

	arr := make([]int, len(src))
	randMax := 0
	idx := 0
	for count := range src {
		randMax += count
		arr[idx] = randMax
		idx++
	}
	// rand := rand.Intn(randMax)
	rand := randMax
	slog.Debug("arr", arr)
	slog.Debug("randMax", randMax)
	slog.Debug("rand", rand)
	idx = sort.Search(len(arr), func(i int) bool {
		return arr[i] >= rand
	})

	slog.Debug("idx", idx)
}

func small() {
	if data, err := fileio.LoadHaikuData("./fileio/small.json"); err != nil {
		panic(err)
	} else {
		corpus := data.CreateCorpus()
		// slog.Debug("corpus", corpus)
		fileio.SaveCorpus("./fileio/small_corpus.json", corpus)
		uta := generator.GenerateFromBegin(corpus)
		slog.Debug("uta", uta)
	}
}

func main() {
	utaFromKigo()
}
