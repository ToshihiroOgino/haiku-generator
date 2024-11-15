package fileio

import (
	"fmt"
	"markov_generator/domain"
	"reflect"
	"testing"
)

func TestLoadHaikuData(t *testing.T) {
	got, err := LoadHaikuData("test.json")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fmt.Printf("got: %v\n", got)

	want := HaikuData{
		"春": {
			"あ": []domain.Uta{
				"あああああいいいいいいいおおおおお",
				"あいうえおかきくけこさしすせそたち",
			},
		},
		"夏": {
			"い": []domain.Uta{
				"いいいい",
				"いいいいいいい",
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
