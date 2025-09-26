/*
Рассмотреть следующий код и ответить на вопросы: к каким негативным последствиям он может привести и как это исправить?
Приведите корректный пример реализации.

var justString string

	func someFunc() {
	  v := createHugeString(1 &lt;&lt; 10)
	  justString = v[:100]
	}

	func main() {
	  someFunc()
	}

Вопрос: что происходит с переменной justString?
*/
package main

import "strings"

var justString string

func someFunc() {
	// v := createHugeString(1 << 10)
	// justString = v[:100]
	/*
		Две потенциальные проблемы есть:
		1. Утечка памяти:
		Так как копирование выполнено не явно, то строка justString и v ссылаются в памяти
		на один и тот же базовый массив, слайс justString является подслайсом v и до тех пор пока
		будет использоваться justString v не будет собран коллектором мусора и будет засорять ОЗУ,
		что является утечкой памяти. Чем больше v тем серъезнее проблема, а если такое действие
		проводится в цикле, то все еще хуже.

		Чтобы в памяти создался новый базовый массив нужно выполнить явное копирование или любую другую операцию,
		которая повлечет за собой создание нового базового массива в памяти.

		2. Возможность выйти за пределы при взятии подслайса v:
		Если при действии v[:n], n > len(v), то получим рантайм ошибку, если поменять n = 2000
		panic: runtime error: slice bounds out of range [:2000] with length 1024
	*/

	v := createHugeString(1 << 10)

	// Явное копирование с проверкой на длинну строки
	justString = safeSlice(v, 0, 100)

	// Конвертация в массив байтов и обратно в строку
	justString = string([]byte(v[:100]))

}

func createHugeString(size int) string {
	return strings.Repeat("x", size)
}

func safeSlice(v string, start, length int) string {
	if start < 0 || length <= 0 || start >= len(v) {
		return ""
	}
	if start+length > len(v) {
		length = len(v) - start
	}

	sl := make([]byte, length)
	copy(sl, v[start:start+length])
	return string(sl)
}

func main() {
	someFunc()
}
