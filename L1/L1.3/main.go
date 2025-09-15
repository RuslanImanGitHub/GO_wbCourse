/*
Реализовать постоянную запись данных в канал (в главной горутине).

Реализовать набор из N воркеров, которые читают данные из этого канала и выводят их в stdout.

Программа должна принимать параметром количество воркеров и при старте создавать указанное число горутин-воркеров.
*/

package main

import(
	"fmt"
	"sync"
	"math/rand"
	"time"
	"os"
	"strconv"
)

// Запуск командой go run main.go N, тут N - количество воркеров запускаемых при старте программы
func main() {
	done := make(chan bool, 1)
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
	}
	go mainGoroutine(n, done)
	<- done
}

func mainGoroutine(n int, done chan bool) {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	
	// Буферинг нужен для асинхронного доступа к одному и тому же ресурсу без дедлоков
	data := make(chan string, 1) 

	for i := 0; i < n; i++ {
		wg.Add(1)
		
		data <- randStr(10)
		go worker(i, data, &wg)
	}
	wg.Wait()
	close(data)
	done <- true
}

func worker(n int, data <- chan string , wg *sync.WaitGroup) {
	fmt.Println("Worker ", n, "working! With data: ", <-data)
	defer wg.Done()
}

func randStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}