/*
Реализовать структуру-счётчик, которая будет инкрементироваться в конкурентной среде (т.е. из нескольких горутин).
По завершению программы структура должна выводить итоговое значение счётчика.

Подсказка: вам понадобится механизм синхронизации, например, sync.Mutex или sync/Atomic для безопасного инкремента.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

func main() {
	var counter atomic.Int64
	var checkVal atomic.Int64
	var wg sync.WaitGroup

	minVal := 10
	maxVal := 20
	rndIntInRange := minVal + rand.Intn(maxVal-minVal+1)
	fmt.Println("rndIntInRange: ", rndIntInRange)
	for i := range rndIntInRange {
		wg.Add(1)
		go func(i int) {
			source := rand.NewSource(rand.Int63())
			rng := rand.New(source)
			rndCntInRange := minVal + rng.Intn(maxVal-minVal+1)
			fmt.Println("rndCntInRange", i, ": ", rndCntInRange)
			for range rndCntInRange {
				counter.Add(1)
			}
			checkVal.Add(int64(rndCntInRange))
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("Counter: ", counter.Load())
	fmt.Println("CheckValue: ", checkVal.Load())
}
