package stats

import (
	. "markov_generator/domain"
	"reflect"
	"testing"
)

func TestCreateVocabularySet(t *testing.T) {
	tests := []struct {
		name      string
		in        []*Morpheme
		want      map[Surface]Vocabulary
		expectErr bool
	}{
		{
			name: "normal",
			in: []*Morpheme{
				{Surface: "塵取", Reading: "チリトリ"},
				{Surface: "に", Reading: "ニ"},
				{Surface: "はこび", Reading: "ハコビ"},
				{Surface: "て", Reading: "テ"},
			},
			want: map[Surface]Vocabulary{
				"塵取":  {ID: 0, Reading: "チリトリ"},
				"に":   {ID: 1, Reading: "ニ"},
				"はこび": {ID: 2, Reading: "ハコビ"},
				"て":   {ID: 3, Reading: "テ"},
			},
			expectErr: false,
		},
		{
			name: "duplicate should return error",
			in: []*Morpheme{
				{Surface: "塵取", Reading: "チリトリ"},
				{Surface: "塵取", Reading: "アアア"},
			},
			want:      nil,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateVocabularySet(tt.in)
			want := tt.want
			if err != nil && !tt.expectErr {
				t.Errorf("CreateVocabularySet() error = %v, wantErr %v", err, tt.expectErr)
			} else if !reflect.DeepEqual(got, want) {
				t.Errorf("CreateVocabularySet() = %v, want %v", got, want)
			}
		})
	}
}
