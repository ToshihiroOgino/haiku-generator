package main

import (
	"fmt"
	"markov_generator/domain"
	"markov_generator/fileio"
	"markov_generator/mecab"
	"regexp"

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

func fn3() {
	str := "「アンタレスの食」てふ過ぎて夜の桜"
	removeBracket, err := regexp.Compile(`[(（][^)）]+[)）]`)
	if err != nil {
		panic(err)
	}
	removeJunk, err := regexp.Compile(`[／　]`)
	if err != nil {
		panic(err)
	}
	str = removeBracket.ReplaceAllString(removeJunk.ReplaceAllString(str, ""), "")
	fmt.Println(str)
}

func fn4() {
	u1 := domain.Uta("帰らなんいざ草の庵は春の風　（学校（教師）をやめる）")
	u2 := domain.Uta("下丸子の30花の天")
	u3 := domain.Uta("下丸子の花の天")
	slog.Debug("u1:", u1.GetCleaned())
	slog.Debug("u2:", u2.IsInvalid())
	slog.Debug("u3:", u3.IsInvalid())
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

func main() {
	corpus()
}
