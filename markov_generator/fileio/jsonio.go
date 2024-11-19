package fileio

import (
	"encoding/json"
	"markov_generator/stats"
	"os"

	"github.com/gookit/slog"
)

func LoadHaikuData(path string) (stats.HaikuData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	slog.Info("reading haiku data from " + path)
	decoder := json.NewDecoder(file)
	var res stats.HaikuData
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	slog.Info("haiku data loaded")
	return res, nil
}

func SaveVocabularySet(path string, vocab stats.VocabularySet) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	slog.Info("saving vocabulary set to " + path)

	jsonEncodable := map[string]stats.VocabID{}
	for k, v := range vocab {
		jsonEncodable[k.String()] = v.ID
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(jsonEncodable)
	if err != nil {
		return err
	}
	slog.Info("vocabulary set saved")
	return nil
}

func SaveCorpus(path string, corpus *stats.Corpus) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	slog.Info("saving corpus to " + path)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(corpus)
	if err != nil {
		return err
	}
	slog.Info("corpus saved")

	return nil
}

func LoadCorpus(path string) (*stats.Corpus, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	slog.Info("reading corpus from " + path)
	decoder := json.NewDecoder(file)
	var res stats.Corpus
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	slog.Info("corpus loaded")
	return &res, nil
}

func SaveKigoStat(path string, kigoStat *stats.KigoStat) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	slog.Info("saving kigo stat to " + path)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(kigoStat)
	if err != nil {
		return err
	}
	slog.Info("kigo stat saved")

	return nil
}
