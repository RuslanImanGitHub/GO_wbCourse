/*
Поменять местами два числа без использования временной переменной.
Подсказка: примените сложение/вычитание или XOR-обмен.
*/

package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	a := rand.Int()
	b := rand.Int()

	fmt.Printf("originals: %d , %d\n", a, b)

	// 1. Встроенный механизм
	a, b = b, a
	fmt.Printf("swapped: %d , %d\n", a, b)
	a, b = b, a
	// 2. Сложение-вычитание (Опасность переполнения)
	arithmeticsSwap(&a, &b)
	fmt.Printf("arithmeticsSwap: %d , %d\n", a, b)
	// 3. XOR-обмен
	xorSwap(&a, &b)
	fmt.Printf("xorSwap: %d , %d\n", a, b)
}

func arithmeticsSwap(a, b *int) {
	if !flowCheck(a, b) {
		fmt.Println("Overflow detected, switching to XOR swap")
		xorSwap(a, b)
		return
	}
	*a = *a + *b
	*b = *a - *b
	*a = *a - *b
}

func xorSwap(a, b *int) {
	*a = *a ^ *b
	*b = *a ^ *b
	*a = *a ^ *b
}

func flowCheck(a, b *int) bool {
	if *b > 0 && *a > math.MaxInt-*b {
		return false // overflow
	}
	if *b < 0 && *a < math.MinInt-*b {
		return false // underflow
	}
	return true
}
