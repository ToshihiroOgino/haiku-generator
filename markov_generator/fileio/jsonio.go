package fileio

import (
	"encoding/json"
	"os"
)

type HaikuData map[string](map[string]([]string))

func LoadHaikuData(path string) (HaikuData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var res HaikuData
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
