/*
Разработать программу, которая проверяет, что все символы в строке
встречаются один раз (т.е. строка состоит из уникальных символов).

Вывод: true, если все символы уникальны, false, если есть повторения.
Проверка должна быть регистронезависимой, т.е. символы в разных регистрах считать одинаковыми.

Например: "abcd" -> true, "abCdefAaf" -> false (повторяются a/A), "aabcd" -> false.

Подумайте, какой структурой данных удобно воспользоваться для проверки условия.
*/
package main

import (
	"fmt"
	"math/rand"
	"unicode"
)

func main() {
	sequence := generateSequence(5, 20)
	fmt.Println("Исходная строка:", sequence)
	unique, nonUniqueChars := AllCharsUnique(sequence)
	fmt.Println("Все символы уникальны?", unique)
	if len(nonUniqueChars) != 0 {
		fmt.Print("Неуникальные символы: ")
		for _, val := range nonUniqueChars {
			fmt.Printf("%c ", val)
		}

	}
}

func AllCharsUnique(str string) (bool, []rune) {
	charMap := make(map[rune]bool)
	var nonUniqueRunes []rune

	for _, char := range str {
		lowerChar := unicode.ToLower(char)

		if charMap[lowerChar] {
			nonUniqueRunes = append(nonUniqueRunes, char)
		}
		charMap[lowerChar] = true
	}

	if len(nonUniqueRunes) != 0 {
		return false, nonUniqueRunes
	} else {
		return true, nil
	}

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
