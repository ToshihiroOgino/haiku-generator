package fileio

import (
	"encoding/json"
	"markov_generator/domain"
	"os"

	"github.com/gookit/slog"
)

type HaikuData map[domain.Season](map[domain.Kigo]([]domain.Uta))

func LoadHaikuData(path string) (HaikuData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	slog.Info("reading haiku data from " + path)
	decoder := json.NewDecoder(file)
	var res HaikuData
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	slog.Info("haiku data loaded")
	return res, nil
}
