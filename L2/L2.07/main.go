/*
Что выведет программа?
Объяснить работу конвейера с использованием select.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

/*
Вывод программы: 1||2 - 3||4- 5||6- 7||8 все числа будут выведены, но заранее нельзя сказать в каком порядке

Что происходит в программе:
asChan() - создает канал и отправляет числа в этот канал с задержками
merge() - объединяет два канала в один с помощью select'а, но дожидаясь пока числа в каналах закончатся
case v, ok := <-a:
	if ok {
		c <- v
	} else {
		a = nil - Закрытие канала, по истечении чисел в канале
	}

Когда оба канала закрыты, то закрывается объединенный канал и его значения выводятся в консоль.

asChan(1,3,5,7) → a channel  ↓
                            select → merged channel → output
asChan(2,4,6,8) → b channel  ↑
*/
