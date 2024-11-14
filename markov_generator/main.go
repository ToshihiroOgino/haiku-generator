package main

import (
	"fmt"
	"markov_generator/mecab"
	"strings"
)

func fn() {
	instance := mecab.CreateInstance()

	haiku := [...]string{"塵取にはこびて藍を植ゑにけり", "明日植うる藍の宵水たつぷりと"}
	for _, h := range haiku {
		res := instance.Parse(h)
		fmt.Printf("input: %s\n%s", h, res)
	}

	instance.Close()
}

func main() {
	line := `明日    名詞,普通名詞,副詞可能,*,*,*,アス,明日,明日,アス,明日,アス,和,*,*,*,*,*,*,体,アス,アス,アス,アス,"2,0",C4,*,184451732742656,671`
	fmt.Println(line)
	arr := strings.Split(line, " ")
	fmt.Println(arr)
	fmt.Println(len(arr))
}
