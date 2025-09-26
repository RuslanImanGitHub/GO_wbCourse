/*
Разработать программу, которая в runtime способна определить тип переменной, переданной в неё (на вход подаётся interface{}).
Типы, которые нужно распознавать: int, string, bool, chan (канал).
Подсказка: оператор типа switch v.(type) поможет в решении.
*/

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num int
	var word string
	var flag bool
	intChan := make(chan int)
	stringChan := make(chan string)
	boolChan := make(chan bool)

	mix := []any{num, word, flag, intChan, stringChan, boolChan}

	for _, val := range mix {
		findType(val)
	}
}

// По опыту C# лучше всего с типами через рефлексию работать, так невозможно опечататься
func findType(obj interface{}) {
	t := reflect.TypeOf(obj).Kind()
	switch t {
	case reflect.Int:
		fmt.Printf("%v is of %s type\n", obj, t)
	case reflect.String:
		fmt.Printf("%v is of %s type\n", obj, t)
	case reflect.Bool:
		fmt.Printf("%v is of %s type\n", obj, t)
	case reflect.Chan:
		eType := reflect.TypeOf(obj).Elem()
		fmt.Printf("%v is of %s %s type\n", obj, t, eType)
	}
}
