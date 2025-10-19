/*
Реализовать функцию, которая будет объединять один или более каналов done (каналов сигнала завершения) в один.
Возвращаемый канал должен закрываться, как только закроется любой из исходных каналов.

Сигнатура функции может быть такой:

Пример использования функции:

	sig := func(after time.Duration) <-chan interface{} {
	   c := make(chan interface{})
	   go func() {
	      defer close(c)
	      time.Sleep(after)
	   }()
	   return c
	}

start := time.Now()
<-or(

	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),

)
fmt.Printf("done after %v", time.Since(start))

В этом примере канал, возвращённый or(...), закроется через ~1 секунду,
потому что самый короткий канал sig(1*time.Second) закроется первым.
Ваша реализация or должна уметь принимать на вход произвольное число
каналов и завершаться при сигнале на любом из них.

Подсказка: используйте select в бесконечном цикле для чтения из всех каналов
одновременно, либо рекурсивно объединяйте каналы попарно.
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(

		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(3*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	// Обработка 0 и 1
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan interface{})

	go func() {
		defer close(done)

		switch len(channels) {
		case 2: //Ожидаем закрытия любого из 2 каналов
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //В случае 3+ каналов рекурсивно дробим на 1-2 канала и ожидаем завершения
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()

	return done
}
