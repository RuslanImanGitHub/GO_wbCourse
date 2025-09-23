/*
	Программа должна корректно завершаться по нажатию Ctrl+C (SIGINT).
	Выберите и обоснуйте способ завершения работы всех горутин-воркеров при получении сигнала прерывания.
	Подсказка: можно использовать контекст (context.Context) или канал для оповещения о завершении.
*/

package main

import (
	"context"
    "os"
    "os/signal"
    "sync"
    "time"
	"fmt"
)

func main() {
	var wg sync.WaitGroup
	// Контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go worker(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Ожидаем вызова SIGINT (Ctrl + C)
	<- sigChan
	cancel()

	fmt.Println("App closing...")
	wg.Wait()
}

// Тикер пишет сообщения в консоль каждую секунду
func worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Проверка завершения контекста
	for {
		select {
		case <- ctx.Done():
			return
		case timer := <- ticker.C:
			fmt.Println("Worker is ticking... ", timer)
		}
	}
}