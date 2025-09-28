/*
Удалить i-ый элемент из слайса. Продемонстрируйте корректное удаление без утечки памяти.
Подсказка: можно сдвинуть хвост слайса на место удаляемого элемента (copy(slice[i:], slice[i+1:])) и уменьшить длину слайса на 1.
*/
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	intSlice := generateSequence(10, 20, 5, 100)
	fmt.Println("original:", intSlice)

	randIndexRemoved := rand.Intn(len(intSlice))
	fmt.Println("remove element", intSlice[randIndexRemoved], "at index", randIndexRemoved)

	result := removeElementAt(intSlice, randIndexRemoved)
	fmt.Println("result:", result)
}

func generateSequence(minElements int, maxElements int, minValue int, maxValue int) (result []int) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, (minValue + rand.Intn(maxValue-minValue+1)))
	}
	return result
}

func removeElementAt(arr []int, i int) []int {
	if i < 0 || i >= len(arr) { //i out of bounds
		return arr
	}

	copy(arr[i:], arr[i+1:]) //Сдвиг элементов

	//Важное действие ↓
	arr[len(arr)-1] = 0

	/*
		slice: [ptrA, ptrB, ptrC] → dataA dataB dataC

		После copy(arr[i:], arr[i+1:])
		slice: [ptrA, ptrC, ptrC] → dataA dataC dataC

		Действие arr[len(arr)-1] = 0
		slice: [ptrA, ptrC, nil] → dataA dataC

		Возвращение arr[:len(arr)-1]
		новый slice: [ptrA, ptrC] → dataA dataC
	*/

	return arr[:len(arr)-1] //Укороченный массив
}
