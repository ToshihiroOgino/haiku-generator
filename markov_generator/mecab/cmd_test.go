package mecab

import (
	"markov_generator/domain"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestParseResult(t *testing.T) {
	tests := []struct {
		line     string
		expected *domain.Morpheme
		hasError bool
	}{
		{
			line:     `塗っ      動詞,一般,*,*,文語四段-ラ行,連用形-促音便,ヌル,塗る,塗っ,ヌッ,塗る,ヌル,和,*,*,*,*,*,*,用,ヌッ,ヌル,ヌッ,ヌル,0,C4,*,7901107804316292,28744`,
			expected: &domain.Morpheme{Surface: "塗っ", Reading: "ヌッ"},
			hasError: false,
		},
		{
			line:     `て        助詞,接続助詞,*,*,*,*,テ,て,て,テ,て,テ,和,*,*,*,*,*,*,接助,テ,テ,テ,テ,*,"動詞%F1,形容詞%F2@-1",*,6837321680953856,24874`,
			expected: &domain.Morpheme{Surface: "て", Reading: "テ"},
			hasError: false,
		},
		{
			line:     `ギヤマン  名詞,普通名詞,一般,*,*,*`,
			expected: &domain.Morpheme{Surface: "ギヤマン", Reading: "ギヤマン"},
			hasError: false,
		},
		{
			line:     `EOF`,
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := parseResult(test.line)
		if test.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}
