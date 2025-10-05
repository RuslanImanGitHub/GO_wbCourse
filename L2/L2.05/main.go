/*
Что выведет программа?
Объяснить вывод программы.
*/
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

/*
Вывод программы: error

Проблема: проверка err != nil всегда будет возвращать true. Даже в IDE выдается предупреждение:
tautological condition: non-nil != nil

Проблема ровно такая же как в L2.03, проверка ==, != смотрит как на тип так и на значение.
test() возвращает ошибку уже с типом customError и проверка != всегда будет возвращать true
*/
