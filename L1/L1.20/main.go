/*
Разработать программу, которая переворачивает порядок слов в строке.
Пример: входная строка:
«snow dog sun», выход: «sun dog snow».
Считайте, что слова разделяются одиночным пробелом. Постарайтесь не использовать дополнительные срезы, а выполнять операцию «на месте».
*/
package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	sequence := generateSequence(10, 20)
	reversedV1 := reverseSentence(sequence)
	reversedV2 := reverseWords(sequence)

	fmt.Println(sequence)
	fmt.Println("Переворот по разбиению пробелов:\n ", reversedV1)
	fmt.Println("Переворот \"in-place\" через манипуляции с рунами:\n ", reversedV2)
}

// То что я сперва написал, просто и понятно разбиватиь строку по пробелам и переставлять слова
// По читаемости мне этот вариант нравится гораздо больше :)
func reverseSentence(str string) string {
	words := strings.Split(str, " ")

	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}

func reverseWords(sentence string) string {
	// Eqmpty or only space string
	if len(strings.TrimSpace(sentence)) == 0 {
		return sentence
	}
	runes := []rune(sentence)

	reverseString(runes, 0, len(runes)-1) //Переворачиваем всю строку

	//Детектим пробелы и переворачиваем строки между ними, чтобы получить оригинальное слово
	//Проблема, что последнее слово не увидим
	counter := 0
	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' {
			reverseString(runes, counter, i-1)
			counter = i + 1
		}
	}
	reverseString(runes, counter, len(runes)-1) //Обработка последнего слова

	return string(runes)
}

// Функция считывает массив символов и переставляет их местами, переворачивая строку
func reverseString(str []rune, start, end int) {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
}

func generateSequence(minElements int, maxElements int) string {
	var result []string
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		animal := gofakeit.Word()

		result = append(result, animal)
	}
	return strings.Join(result, " ")
}
