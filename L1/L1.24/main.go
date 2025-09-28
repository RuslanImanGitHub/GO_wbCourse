/*
Разработать программу нахождения расстояния между двумя точками на плоскости.
Точки представлены в виде структуры Point с инкапсулированными (приватными) полями x, y (типа float64)
и конструктором. Расстояние рассчитывается по формуле между координатами двух точек.

Подсказка: используйте функцию-конструктор NewPoint(x, y), Point и метод Distance(other Point) float64.
*/
package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	minV := -10.00
	maxV := 10.00
	p1 := NewPoint(randFloatInRange(minV, maxV), randFloatInRange(minV, maxV))
	p2 := NewPoint(randFloatInRange(minV, maxV), randFloatInRange(minV, maxV))

	fmt.Println("Point1:", p1)
	fmt.Println("Point2:", p2)
	fmt.Println("Distance 1-2:", p1.Distance(*p2))
	fmt.Println("Distance 2-1:", p2.Distance(*p1))
}

type Point struct {
	x float64
	y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func (p1 *Point) Distance(p2 Point) float64 {
	return math.Sqrt(math.Pow(p2.x-p1.x, 2) + math.Pow(p2.y-p1.y, 2))
}

func randFloatInRange(minValue float64, maxValue float64) float64 {
	return minValue + rand.Float64()*(maxValue-minValue)
}
