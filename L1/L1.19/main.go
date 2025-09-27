/*
Разработать программу, которая переворачивает подаваемую на вход строку.
Например: при вводе строки «главрыба» вывод должен быть «абырвалг».

Учтите, что символы могут быть в Unicode (русские буквы, emoji и пр.),
то есть просто iterating по байтам может не подойти — нужен срез рун ([]rune).
*/
package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func main() {
	sequence := generateSequence(10, 30)
	reversed := reverseString(sequence)

	fmt.Println(sequence)
	fmt.Println(reversed)
}

func reverseString(str string) (result string) {
	var sb strings.Builder
	runes := []rune(str)
	sb.Grow(len(runes))

	for i := len(runes) - 1; i >= 0; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}

func generateSequence(minElements, maxElements int) string {
	const unicodeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]|;:'\",.<>/?`~🥺😂🥰😊😍😝🤗"
	ancientTexts := []rune(unicodeChars)
	var result []rune
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, ancientTexts[rand.Intn(len(ancientTexts)-1)])
	}
	return string(result)
}
