/*
Реализовать собственную функцию sleep(duration) аналогично встроенной функции time.Sleep,
которая приостанавливает выполнение текущей горутины.
Важно: в отличии от настоящей time.Sleep, ваша функция должна именно блокировать выполнение
(например, через таймер или цикл), а не просто вызывать time.Sleep :) — это упражнение.

Можно использовать канал + горутину, или цикл на проверку времени (не лучший способ, но для обучения).
*/
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	lag := 3 * time.Second
	fmt.Println("Время простоя:", lag)

	fmt.Println("Канал с таймером начался")
	CustomSleepChan(lag)
	fmt.Println("Канал с таймером закончился")
	fmt.Println()
	fmt.Println("Упрощенный канал с time.After начался")
	CustomSleepTimer(lag)
	fmt.Println("Упрощенный канал с time.After закончился")
	fmt.Println()
	fmt.Println("Контекст с таймаутом начался")
	CustomSleepContext(lag)
	fmt.Println("Контекст с таймаутом закончился")
}

func CustomSleepChan(duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C //Блокается пока не придет сигнал
}

func CustomSleepTimer(duration time.Duration) {
	<-time.After(duration)
}

func CustomSleepContext(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	<-ctx.Done()
}
