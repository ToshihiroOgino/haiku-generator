package fileio

import (
	"encoding/json"
	"markov_generator/domain"
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

	type data struct {
		Surface domain.Surface
		Reading domain.Reading
		Next    map[stats.VocabID]stats.Count
		Prev    map[stats.VocabID]stats.Count
	}

	corpusEncodable := map[stats.VocabID]data{}
	for k, v := range corpus.VocabularySet {
		corpusEncodable[v.ID] = data{
			Surface: k.Surface,
			Reading: k.Reading,
			Next:    v.Next,
			Prev:    v.Prev,
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(corpusEncodable)
	if err != nil {
		return err
	}
	slog.Info("corpus saved")

	return nil
}
