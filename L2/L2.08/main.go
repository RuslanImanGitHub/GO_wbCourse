package main

import (
	"fmt"
	"ntpTime/ntpTime"
	"os"
	"time"
)

func main() {
	currTime, err := ntpTime.GetCurrentTime()
	if err != nil { // Выводим ошибку в STDERR и выходим с ненулевым кодом
		fmt.Fprintf(os.Stderr, "Ошибка получения времени: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Точное время (NTP): %s\n", currTime.Format(time.RFC3339))
}
