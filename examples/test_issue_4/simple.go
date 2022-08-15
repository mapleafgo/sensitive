package main

import (
	"fmt"

	"mapleafgo.cn/sensitive"
)

func main() {
	filter := sensitive.New()
	filter.LoadWordDict("../../dict/dict.txt")
	fmt.Println(filter.Replace("xC4x", '*'))
}
