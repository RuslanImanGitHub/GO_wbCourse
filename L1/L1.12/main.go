/*
Имеется последовательность строк: ("cat", "cat", "dog", "cat", "tree"). Создать для неё собственное множество.
Ожидается: получить набор уникальных слов. Для примера, множество = {"cat", "dog", "tree"}.
*/

package main

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	//sequence := []string{"cat", "cat", "dog", "cat", "tree"}

	sequence := generateSequence(5, 10)

	uniqueSeq := unique(sequence)

	fmt.Println(sequence)
	fmt.Println(uniqueSeq)
}

func unique(seq []string) (result []string) {
	dict := make(map[string]bool)
	for _, val := range seq {
		dict[val] = true
	}
	for key := range dict {
		result = append(result, key)
	}
	return result
}

func generateSequence(minElements int, maxElements int) (result []string) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		animal := gofakeit.Animal()
		if rand.Intn(2) == 1 { //  random bool
			result = append(result, animal)
		}
		result = append(result, animal)
	}
	return result
}
