/*
Реализовать безопасную для конкуренции запись данных в структуру map.
Подсказка: необходимость использования синхронизации (например, sync.Mutex или встроенная concurrent-map
!!![concurrent-map не встроенная библиотека https://github.com/orcaman/concurrent-map. В GO 1.9 добавилась встроенная библиотека называется sync.Map]).
Проверьте работу кода на гонки (util go run -race).
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Mutex method:")
	MutexMethod()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("SyncMap method:")
	syncMapMethod()
}

// MUTEX PART
type AsyncMap[K comparable, V any] struct {
	mu   sync.RWMutex // Отдельные локи на чтение и запись, позволяет читать многим сразу или писать одному
	data map[K]V
}

func NewAsyncMap[K comparable, V any]() *AsyncMap[K, V] {
	return &AsyncMap[K, V]{
		data: make(map[K]V),
	}
}

// Безопасное чтение
func (am *AsyncMap[K, V]) Get(key K) (V, bool) {
	am.mu.RLock()
	defer am.mu.RLocker().Unlock()
	val, exist := am.data[key]
	return val, exist
}

// Безопасная запись
func (am *AsyncMap[K, V]) Set(key K, value V) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.data[key] = value
}

func MutexMethod() {
	aMap := NewAsyncMap[int, string]()
	var wg sync.WaitGroup

	//Async Writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			val := fmt.Sprintf("Value of key %d", n)
			aMap.Set(n, val)
		}(i)
	}
	wg.Wait()

	//Async Reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if value, exists := aMap.Get(n); exists {
				fmt.Println(n, value)
			}
		}(i)
	}
	wg.Wait()
}

//----------------------------------------------------------

// SYNC.MAP PART
func syncMapMethod() {
	var sMap sync.Map
	var wg sync.WaitGroup

	//Async Writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			val := fmt.Sprintf("Value of key %d", n)
			sMap.Store(n, val)
		}(i)
	}
	wg.Wait()

	//Async READS (Делается снапшот по которому будет идти итерация)
	sMap.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}
