/*
Напишите функцию, которая находит все множества анаграмм по заданному словарю.

Требования
На вход подается срез строк (слов на русском языке в Unicode).
На выходе: map-множество -> список, где ключом является первое встреченное слово множества,
а значением — срез из всех слов, принадлежащих этому множеству анаграмм, отсортированных по возрастанию.

Множества из одного слова не должны выводиться (т.е. если нет анаграмм, слово игнорируется).
Все слова нужно привести к нижнему регистру.

Пример:
Вход: ["пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"]
Результат (ключи в примере могут быть в другом порядке):
– "пятак": ["пятак", "пятка", "тяпка"]
– "листок": ["листок", "слиток", "столик"]

Слово «стол» отсутствует в результатах, так как не имеет анаграмм.

Для решения задачи потребуется умение работать со строками, сортировать
и использовать структуры данных (map).

Оценим эффективность: решение должно работать за линейно-логарифмическое время относительно
количества слов (допустимо n * m log m, где m — средняя длина слова для сортировки букв).
*/
package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func main() {
	//sequence := generateSequence(10, 20)
	sequence := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(sequence)

	anagrams := FindAnagrams(sequence)
	for key, val := range anagrams {
		fmt.Printf("%s: %s\n", key, val)
	}
	fmt.Println()
	anagrams2 := FindAnagramsSort(sequence)
	for key, val := range anagrams2 {
		fmt.Printf("%s: %s\n", key, val)
	}
}

/*func generateSequence(minElements, maxElements int) []string {
	//TODO: gofakeit - generate sequence + force anagrams
}*/

// #region Letter sorting, complexity O(n * m log m)
func FindAnagramsSort(words []string) map[string][]string {
	if len(words) == 0 {
		return make(map[string][]string)
	}
	lowercaseWords := make([]string, len(words))
	for i, val := range words {
		lowercaseWords[i] = strings.ToLower(val)
	}
	anagramMap := make(map[string][]string)
	for _, word := range lowercaseWords {
		sortedWord := sortLetters(word)
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}
	result := make(map[string][]string)
	for _, values := range anagramMap {
		if len(values) <= 1 {
			continue
		}

		sort.Strings(values)
		result[values[0]] = values
	}
	return result
}

func sortLetters(word string) string {
	runes := []rune(word)
	sort.Slice(runes, func(i int, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

// #endregion

// #region Word comparison, complexity bad - O(n^2 * m)
func FindAnagrams(words []string) map[string][]string {
	result := make(map[string][]string)

	if len(words) == 0 {
		return result
	}
	lengthSlice := make([]int, 0)
	lowercaseWords := make([]string, len(words))
	for i, val := range words {
		lowercaseWords[i] = strings.ToLower(val)
		if !slices.Contains(lengthSlice, len(val)) {
			lengthSlice = append(lengthSlice, len(val))
		}
	}

	groupedWords := getWordsByLength(lowercaseWords, lengthSlice)
	for _, words := range groupedWords {
		processed := make(map[string]bool)
		for i, word1 := range words {
			if processed[word1] {
				continue
			}

			var anagramGroup []string
			anagramGroup = append(anagramGroup, word1)
			processed[word1] = true

			for j := i + 1; j < len(words); j++ {
				word2 := words[j]
				if processed[word2] {
					continue
				}

				if areAnagrams(word1, word2) {
					anagramGroup = append(anagramGroup, word2)
					processed[word2] = true
				}
			}

			if len(anagramGroup) > 1 {
				sort.Strings(anagramGroup)
				result[anagramGroup[0]] = anagramGroup
			}
		}
	}

	return result
}

func getWordsByLength(words []string, lengths []int) map[int][]string {
	result := make(map[int][]string, len(lengths))

	for _, length := range lengths {
		for _, word := range words {
			if len(word) == length {
				result[length] = append(result[length], word)
			}
		}
	}

	return result
}

func areAnagrams(word1, word2 string) bool {
	map1 := make(map[rune]int)
	map2 := make(map[rune]int)

	for _, letter := range word1 {
		map1[letter]++
	}
	for _, letter := range word2 {
		map2[letter]++
	}

	for letter, count1 := range map1 {
		if count2, ok := map2[letter]; !ok || count1 != count2 {
			return false
		}
	}
	for letter, count2 := range map2 {
		if count1, ok := map1[letter]; !ok || count2 != count1 {
			return false
		}
	}

	return true
}

// #endregion
