package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"regexp"
)

func main() {
	// 从剪贴板读取数据
	clipboardData, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("无法读取剪贴板数据:", err)
		return
	}

	// 使用正则表达式拆分数据并输出
	wordPattern := regexp.MustCompile(`\s+|(?<=[^\p{L}])|(?=[^\p{L}])`)
	words := wordPattern.Split(clipboardData, -1)
	for _, word := range words {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
}
