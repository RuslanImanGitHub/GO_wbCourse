/*
Разработать программу, которая перемножает, делит, складывает,
вычитает две числовых переменных a, b, значения которых > 2^20 (больше 1 миллион).

Комментарий: в Go тип int справится с такими числами, но обратите внимание на
возможное переполнение для ещё больших значений. Для очень больших чисел можно использовать math/big.
*/
package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

func main() {
	randBigNumber1 := generateBigInt()
	randBigNumber2 := generateBigInt()

	fmt.Println("a = ", randBigNumber1)
	fmt.Println("b = ", randBigNumber2)

	operations := new(big.Float)

	operations.Add(randBigNumber1, randBigNumber2)
	fmt.Println("a + b = ", operations.String())

	operations.Sub(randBigNumber1, randBigNumber2)
	fmt.Println("a - b = ", operations.String())

	operations.Mul(randBigNumber1, randBigNumber2)
	fmt.Println("a * b = ", operations.String())

	operations.Quo(randBigNumber1, randBigNumber2)
	fmt.Println("a / b = ", operations.String())
}

func generateBigInt() *big.Float {
	min := int64(1 << 40)
	max := int64(1 << 50)

	return big.NewFloat(float64(min + rand.Int63n(max-min+1)))
}
