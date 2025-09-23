/*
Дана структура Human (с произвольным набором полей и методов).

Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).

Подсказка: используйте композицию (embedded struct), чтобы Action имел все методы Human.
*/

package main

import "fmt"

type Action struct {
	actionDescription string 
}

type Human struct {
	name string
	age int
	Action
}

func NewHuman(name string, age int, desc string) *Human {
	human := Human{name: name, age: age, Action: Action{actionDescription: desc}}
	return &human
}

func (h *Human) Print() {
	fmt.Printf("Human: %s, %d, %s\n", h.name, h.age, h.Action.actionDescription)
}

func main() {
	var humanIvan = NewHuman("Ivan", 23, "Work")
	humanIvan.Print()

	var humanJohn = NewHuman("Jhon", 7, "Play")
	humanJohn.Print()

	var humanDasha = NewHuman("Dasha", 14, "Learn")
	humanDasha.Print()
}