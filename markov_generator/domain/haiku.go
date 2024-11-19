package domain

import "regexp"

type Uta string
type Kigo string
type Season string

func (u Uta) String() string {
	return string(u)
}

func (k Kigo) String() string {
	return string(k)
}

var invalidSymbols = regexp.MustCompile(`[「『＼≪\?？！…“♯♭、〇〈0-9０-９a-zA-Zａ-ｚＡ-Ｚ～ﾐ＝・]`)

func (u Uta) IsInvalid() bool {
	return invalidSymbols.MatchString(u.String())
}

var bracket = regexp.MustCompile(`[(（][^)）]+[)）]`)
var removableSymbols = regexp.MustCompile(`[／　/\*]`)

func (u Uta) GetCleaned() Uta {
	return Uta(bracket.ReplaceAllString(removableSymbols.ReplaceAllString(u.String(), ""), ""))
}
