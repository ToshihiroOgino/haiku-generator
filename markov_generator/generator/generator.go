package generator

import (
	"markov_generator/domain"
	"markov_generator/mecab"
	"markov_generator/stats"
	"math/rand"
	"sort"

	"github.com/gookit/slog"
)

type GeneratedUta struct {
	Uta       domain.Uta
	Reading   string
	UtaLength int
}

func GenerateFromBegin(corpus *stats.Corpus) GeneratedUta {
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
		length += m.Reading.Length()
		utaReading += m.Reading.String()
		prevID = nextID
		// slog.Debug("uta", uta)
	}
	return GeneratedUta{domain.Uta(uta), utaReading, length}
}

func GenerateFromKigo(c *stats.Corpus, kigoStat *stats.KigoStat, kigo domain.Kigo, instance *mecab.MeCab) GeneratedUta {
	// 季語を埋め込む位置の候補を探索
	averageKigoPos := -1
	for _, kigoInfoArr := range kigoStat.Individual {
		for _, info := range kigoInfoArr {
			if info.Kigo == kigo {
				averageKigoPos = int(info.AveragePos)
				break
			}
		}
		if averageKigoPos != -1 {
			break
		}
	}

	utaMorphemeArr := instance.Exec(string(kigo))
	kigoLength := 0
	for _, m := range utaMorphemeArr {
		kigoLength += m.Reading.Length()
	}

	posCandidate := [4]int{0, 5, 12, averageKigoPos}
	r := rand.Intn(len(posCandidate))
	kigoPos := posCandidate[r]

	// 5+7+5-kigoLength は季語が17音の俳句に埋め込まれる場合の最大の文字数
	if kigoPos < -1 || kigoPos >= 5+7+5-kigoLength {
		// 候補1 ランダム
		// 12は俳句の5_7_5の前半2つ(5_7)と最後の節の頭までの間に俳句を生成するための最大の文字数
		// kigoPos = rand.Intn(12)

		// 候補2 全体平均±誤差2
		// 4はデータセット全体の平均位置(kigoStat.TotalAverage)の整数部分
		// kigoPos = 4 + rand.Intn(5) - 2

		// 候補3 各節の冒頭のいづれか
		posCandidate := [3]int{0, 5, 12}
		r := rand.Intn(len(posCandidate))
		kigoPos = posCandidate[r]
	}

	// 季語の前後から俳句を生成
	const MAX_RETRY = 10
	// 季語より前の句
	wantLength := kigoPos
	nextID := c.GetIDFromMorpheme(&utaMorphemeArr[0])
	for wantLength > 0 {
		// いくつか候補を選び、字余りとならない形態素を選択する
		// 超過しない候補
		okCandidates := []*domain.Morpheme{}
		// 超過する候補
		ngCandidates := []*domain.Morpheme{}
		for retryCount := 0; retryCount < MAX_RETRY; retryCount++ {
			id := generatePrev(c, nextID)
			b := c.GetMorphemeFromID(id)
			if id == c.Eos.ID || id == c.Bos.ID {
				continue
			} else if wantLength-b.Reading.Length() >= 0 {
				okCandidates = append(okCandidates, b)
			} else {
				ngCandidates = append(ngCandidates, b)
			}
		}

		var m *domain.Morpheme
		if len(okCandidates) == 0 && len(ngCandidates) == 0 {
			break
		} else if len(okCandidates) > 0 {
			// 超過しない候補があればそれを選択
			idx := rand.Intn(len(okCandidates))
			m = okCandidates[idx]
		} else {
			// 超過する候補のみならら、その中からランダムに選択
			idx := rand.Intn(len(ngCandidates))
			m = ngCandidates[idx]
		}

		wantLength -= m.Reading.Length()
		// 先頭に追加
		utaMorphemeArr = append([]domain.Morpheme{*m}, utaMorphemeArr...)
		nextID = c.GetIDFromMorpheme(m)
	}
	// 季語より後の句
	// すでに出来上がっている句を字余りにしないような長さで後の句を生成したい
	prevID := c.GetIDFromMorpheme(&utaMorphemeArr[len(utaMorphemeArr)-1])
	wantLength = (5 + 7 + 5) - (kigoPos + kigoLength) + wantLength
	for wantLength > 0 {
		okCandidates := []*domain.Morpheme{}
		ngCandidates := []*domain.Morpheme{}
		for retryCount := 0; retryCount < MAX_RETRY; retryCount++ {
			id := generateNext(c, prevID)
			b := c.GetMorphemeFromID(id)
			if id == c.Eos.ID || id == c.Bos.ID {
				continue
			} else if wantLength-b.Reading.Length() >= 0 {
				okCandidates = append(okCandidates, b)
			} else {
				ngCandidates = append(ngCandidates, b)
			}
		}

		var m *domain.Morpheme
		if len(okCandidates) == 0 && len(ngCandidates) == 0 {
			break
		} else if len(okCandidates) > 0 {
			// 超過しない候補があればそれを選択
			idx := rand.Intn(len(okCandidates))
			m = okCandidates[idx]
		} else {
			// 超過する候補のみならら、その中からランダムに選択
			idx := rand.Intn(len(ngCandidates))
			m = ngCandidates[idx]
		}

		wantLength -= m.Reading.Length()
		// 末尾に追加
		utaMorphemeArr = append(utaMorphemeArr, *m)
		prevID = c.GetIDFromMorpheme(m)
	}

	uta := ""
	reading := ""
	utaLength := 0
	for _, m := range utaMorphemeArr {
		uta += m.Surface.String()
		reading += m.Reading.String()
		utaLength += m.Reading.Length()
	}
	return GeneratedUta{domain.Uta(uta), reading, utaLength}
}

func generateNext(c *stats.Corpus, prevID stats.VocabID) stats.VocabID {
	prev := c.GetVocabularyFromID(prevID)
	if prev == nil {
		slog.Warn("prev is nil", "prevID", prevID)
		return c.Eos.ID
	}
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

func generatePrev(c *stats.Corpus, nextID stats.VocabID) stats.VocabID {
	next := c.GetVocabularyFromID(nextID)
	if next == nil {
		slog.Warn("next is nil", "nextID", nextID)
		return c.Bos.ID
	}
	arr := make([]int, len(next.Prev))
	randMax := 0
	idx := 0
	for _, count := range next.Prev {
		randMax += count
		arr[idx] = randMax
		idx++
	}
	if randMax == 0 {
		return c.Bos.ID
	}
	rand := rand.Intn(randMax)
	idx = sort.Search(len(arr), func(i int) bool {
		return arr[i] >= rand
	})
	keys := make([]stats.VocabID, 0, len(next.Prev))
	for vi := range next.Prev {
		keys = append(keys, vi)
	}
	return stats.VocabID(keys[idx])
}
