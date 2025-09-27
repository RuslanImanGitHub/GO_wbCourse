/*
Реализовать алгоритм бинарного поиска встроенными методами языка.
Функция должна принимать отсортированный слайс и искомый элемент,
возвращать индекс элемента или -1, если элемент не найден.

Подсказка: можно реализовать рекурсивно или итеративно, используя цикл for.
*/
package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	upperBound := 10

	sequence := generateSequence(10, 20, 0, upperBound)
	searchVal := rand.Intn(upperBound)
	foundIndex := binarySearchIndex(sequence, searchVal, 0, len(sequence)-1)

	fmt.Println(sequence)
	fmt.Println(searchVal)
	if foundIndex != -1 {
		fmt.Printf("Index - %d, Value - %d\n", foundIndex, sequence[foundIndex])
	} else {
		fmt.Println(foundIndex)
	}

}

func binarySearchIndex(slice []int, target, start, end int) int {
	if start > end {
		return -1
	}
	mid := start + (end-start)/2

	switch {
	case slice[mid] == target:
		return mid
	case slice[mid] < target:
		return binarySearchIndex(slice, target, mid+1, end)
	case slice[mid] > target:
		return binarySearchIndex(slice, target, start, mid-1)
	}
	return -1
}

func generateSequence(minElements, maxElements, minValue, maxValue int) (result []int) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, (minValue + rand.Intn(maxValue-minValue+1)))
	}
	sort.Ints(result)
	return result
}
