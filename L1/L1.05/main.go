/*
	Разработать программу, которая будет последовательно отправлять значения в канал,
	а с другой стороны канала – читать эти значения. По истечении N секунд программа должна завершаться.

	Подсказка: используйте time.After или таймер для ограничения времени работы.
*/

package main

import (
	"fmt"
	"time"
	"strconv"
	"os"
)

// Запуск командой go run main.go N, тут N - через сколько секунд наступит таймаут
func main() {
	// В качестве количества времени возьмем птрибут запуска программы
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	counterChan := make(chan int)

	// Источник отправки значений
	go func() {
		for i := 1; i <= 100 ; i++ {
			counterChan <- i
			time.Sleep(time.Second)
		}
		close(counterChan)
	}()

	timout := time.After(time.Duration(n) * time.Second)

	// Приемник значений
	for {
		select {
		case val, ok := <- counterChan:
			if !ok { // Канал закрыт, прошло более 100 секунд
				break
			}
			fmt.Println("Прием данных: ", val)
		case stopTime := <- timout:
			fmt.Println("Таймаут! Вышло время - ", stopTime)
			return
		}
	}
}