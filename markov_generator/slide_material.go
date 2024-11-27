//go:build ignore

package main

import (
	"encoding/json"
	"markov_generator/domain"
	"markov_generator/fileio"
	"markov_generator/stats"
	"sort"
	"strings"

	"github.com/gookit/slog"
)

func getTop3Keys(m map[stats.VocabID]int) []stats.VocabID {
	type kv struct {
		Key   stats.VocabID
		Value int
	}
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	keys := []stats.VocabID{}
	for i := 0; i < 3; i++ {
		keys = append(keys, ss[i].Key)
	}
	return keys
}

func kigoPrevNext() {
	c, err := fileio.LoadCorpus("./fileio/corpus.json")
	if err != nil {
		slog.FatalErr(err)
	}
	m := &domain.Morpheme{
		Surface: "花野",
		Reading: "ハナノ",
	}
	id := c.GetIDFromMorpheme(m)
	if id == -1 {
		slog.Fatal("Not found morpheme", m)
	}
	info := c.GetVocabularyFromID(id)
	slog.Debug("ID", id)

	nextTop3, prevTop3 := getTop3Keys(info.Next), getTop3Keys(info.Prev)
	slog.Info("nextTop3", nextTop3)
	sumNext := 0
	for _, count := range info.Next {
		sumNext += count
	}
	for _, id := range nextTop3 {
		count := info.Next[id]
		morpheme := c.GetMorphemeFromID(id)
		slog.Debugf("%s %.01f%%", morpheme.Surface, float64(count)/float64(sumNext)*float64(100))
	}

	slog.Info("prevTop3", prevTop3)
	sumPrev := 0
	for _, count := range info.Prev {
		sumPrev += count
	}
	for _, id := range prevTop3 {
		count := info.Prev[id]
		morpheme := c.GetMorphemeFromID(id)
		slog.Debugf("%s %.01f%%", morpheme.Surface, float64(count)/float64(sumPrev)*float64(100))
	}
}

func topKigo() {
	kigoStats, _ := fileio.LoadKigoStat("./fileio/kigo_stat.json")

	allInfo := []*stats.KigoInfo{}
	for _, infoList := range kigoStats.Individual {
		allInfo = append(allInfo, infoList...)
	}
	slog.Debug("len(allInfo)", len(allInfo))
	sort.Slice(allInfo, func(i, j int) bool {
		return allInfo[i].NumHaiku >= allInfo[j].NumHaiku
	})
	bytes, _ := json.MarshalIndent(allInfo[:10], "", "  ")
	slog.Debug("top", string(bytes))
}

func searchMorpheme(c *stats.Corpus, surface domain.Surface) []*domain.Morpheme {
	res := []*domain.Morpheme{}
	for _, m := range c.MorphemeList {
		if m.Surface == surface {
			res = append(res, m)
		}
	}

	for _, m := range searchMorpheme(c, domain.Surface("藤")) {
		slog.Debug(m)
	}

	return res
}

func dataSize() {
	haikuData, _ := fileio.LoadHaikuData("./fileio/haiku.json")
	numKigo := 0
	numHaiku := 0
	maxHaikuPerKigo := 0
	for _, seasonMap := range haikuData {
		for kigo, kigoMap := range seasonMap {
			if len(kigoMap) > 0 {
				numKigo++
			}
			numHaiku += len(kigoMap)
			// 季語と同値な文字列を含む俳句を数える
			count := 0
			for _, haiku := range kigoMap {
				if strings.Index(haiku.String(), kigo.String()) != -1 {
					count++
				}
			}
			maxHaikuPerKigo = max(maxHaikuPerKigo, count)
		}
	}
	slog.Info("numKigo", numKigo, "numHaiku", numHaiku, "maxHaikuPerKigo", maxHaikuPerKigo)
}

func main() {
	// topKigo()
	// kigoPrevNext()
	dataSize()
}
