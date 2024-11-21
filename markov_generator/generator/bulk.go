package generator

import (
	"markov_generator/domain"
	"markov_generator/mecab"
	"markov_generator/stats"
	"math/rand"
	"sort"

	"github.com/gookit/slog"
)

type GeneratedHaikuData map[string](map[domain.Kigo]([]domain.Uta))

func GenerateALot(corpus *stats.Corpus, kigoStat *stats.KigoStat) GeneratedHaikuData {
	rand.NewSource(0)

	allInfo := []*stats.KigoInfo{}
	for _, infoList := range kigoStat.Individual {
		allInfo = append(allInfo, infoList...)
	}
	sort.Slice(allInfo, func(i, j int) bool {
		return allInfo[i].NumHaiku >= allInfo[j].NumHaiku
	})
	const NUM_KIGO = 20
	const NUM_HAIKU = 10

	generated := make(GeneratedHaikuData)
	const DECENT_CATEGORY = "Decent"
	const RANDOM_CATEGORY = "Random"
	generated[DECENT_CATEGORY] = make(map[domain.Kigo]([]domain.Uta))
	generated[RANDOM_CATEGORY] = make(map[domain.Kigo]([]domain.Uta))

	instance := mecab.CreateInstance()
	defer instance.Close()

	slog.Info("generating decent kigo haiku")
	for i := 0; i < NUM_KIGO; i++ {
		kigo := allInfo[i].Kigo
		generated[DECENT_CATEGORY][kigo] = []domain.Uta{}
		for j := 0; j < NUM_HAIKU; j++ {
			uta := GenerateFromKigo(corpus, kigoStat, kigo, instance)
			generated[DECENT_CATEGORY][kigo] = append(generated[DECENT_CATEGORY][kigo], uta.Uta)
		}
	}

	slog.Info("generating random kigo haiku")
	for i := 0; i < NUM_KIGO; i++ {
		randIndex := rand.Intn(len(allInfo))
		kigo := allInfo[randIndex].Kigo
		generated[RANDOM_CATEGORY][kigo] = []domain.Uta{}
		for j := 0; j < NUM_HAIKU; j++ {
			uta := GenerateFromKigo(corpus, kigoStat, kigo, instance)
			generated[RANDOM_CATEGORY][kigo] = append(generated[RANDOM_CATEGORY][kigo], uta.Uta)
		}
	}
	return generated
}
