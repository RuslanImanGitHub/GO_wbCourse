/*
Реализовать алгоритм быстрой сортировки массива встроенными средствами языка.
Можно использовать рекурсию.

Подсказка: напишите функцию quickSort([]int) []int которая сортирует срез целых чисел.
Для выбора опорного элемента можно взять середину или первый элемент.
*/
package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	sequence := generateSequence(5, 20, 0, 100)
	sorted := quickSort(sequence)

	fmt.Println(sequence)
	fmt.Println(sorted)
}

func quickSort(slice []int) (result []int) {
	if len(slice) <= 1 {
		return slice
	}

	pivot := slice[int(math.Floor(float64(len(slice)/2)))]
	var lesser []int
	var bigger []int
	var equal []int

	for _, val := range slice {
		switch {
		case val < pivot:
			lesser = append(lesser, val)
		case val > pivot:
			bigger = append(bigger, val)
		case val == pivot:
			equal = append(equal, val)
		}
	}

	return append(append(quickSort(lesser), equal...), quickSort(bigger)...)
}

func generateSequence(minElements, maxElements, minValue, maxValue int) (result []int) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, (minValue + rand.Intn(maxValue-minValue+1)))
	}
	return result
}
