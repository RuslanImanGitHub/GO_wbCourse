/*
Написать программу, которая конкурентно рассчитает значения квадратов чисел, взятых из массива [2,4,6,8,10], и выведет результаты в stdout.

Подсказка: запусти несколько горутин, каждая из которых возводит число в квадрат.
*/
package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"L1.2/stopwatch"
)

func main() {

	sw := stopwatch.NewStopwatch()
	swas := stopwatch.NewStopwatch()
	arr := []int{2, 4, 6, 8, 10}
	//Sync
	sw.Start()
	for _, value := range arr {
		os.Stdout.WriteString(strconv.Itoa(value*value) + "\n")
	}
	sw.Stop()

	fmt.Printf("Sync pass. ElapsedTime: %s\n", sw.Elapsed())

	//Async
	wg := new(sync.WaitGroup)
	swas.Start()
	for _, value := range arr {
		wg.Add(1)
		square(value, wg)
	}
	swas.Stop()

	fmt.Printf("Async pass. ElapsedTime: %s\n", swas.Elapsed())
}

func square(value int, wg *sync.WaitGroup) {
	os.Stdout.WriteString(strconv.Itoa(value*value) + "\n")
	defer wg.Done()
}
