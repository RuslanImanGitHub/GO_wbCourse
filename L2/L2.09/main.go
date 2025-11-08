/*
Написать функцию Go, осуществляющую примитивную распаковку строки, содержащей повторяющиеся символы/руны.

Примеры работы функции:

Вход: "a4bc2d5e"
Выход: "aaaabccddddde"

Вход: "abcd"
Выход: "abcd" (нет цифр — ничего не меняется)

Вход: "45"
Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)

Вход: ""
Выход: "" (пустая строка -> пустая строка)

Дополнительное задание
Поддерживать escape-последовательности вида \:

Вход: "qwe\4\5"
Выход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)

Вход: "qwe\45"
Выход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)

Требования к реализации
Функция должна корректно обрабатывать ошибочные случаи (возвращать ошибку, например, через error), и проходить unit-тесты.

Код должен быть статически анализируем (vet, golint).
*/

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

func main() {
	str := generateSequence(10, 20)
	fmt.Println(str)
	unpacked, err := UnpackSequenceV2(str)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(unpacked)
}

// Не может собрать цифру из нескольких рун ("a45" → aaaaa) Depricated оставил как часть истории
/*func unpackSequence(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	alphabetOnly := regexp.MustCompile(`^[a-zA-Z]+$`)
	digitOnly := regexp.MustCompile(`^\d+$`)

	if digitOnly.MatchString(str) {
		return "", errors.New("String contains only digits. String must contain either escape sequence (`\\`) with digits or alphanumeric chars")
	}
	if alphabetOnly.MatchString(str) {
		return str, nil
	}

	var sb strings.Builder
	sb.Grow(len(str))
	strRunes := []rune(str)
	var prevRune rune
	for _, v := range strRunes {
		if v == '\\' && prevRune != '\\' {
			prevRune = v
			continue
		}
		if v == '\\' && prevRune == '\\' {
			sb.WriteRune(v)
			prevRune = 0
			continue
		}
		if unicode.IsDigit(v) && prevRune == '\\' {
			sb.WriteRune(v)
			prevRune = v
			continue
		}
		if unicode.IsLetter(v) {
			sb.WriteRune(v)
			prevRune = v
			continue
		}
		if unicode.IsDigit(v) && prevRune != '\\' {
			sb.Grow(int(v-'0') - 1)
			for range int(v-'0') - 1 {
				sb.WriteRune(prevRune)
			}
			prevRune = 0
			continue
		}
	}
	return sb.String(), nil
}*/

// Улучшенаая версия ("a45" → aaaaaaaa...)
func UnpackSequenceV2(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	alphabetOnly := regexp.MustCompile(`^[a-zA-Z]+$`)
	digitOnly := regexp.MustCompile(`^\d+$`)

	if digitOnly.MatchString(str) {
		return "", errors.New("string contains only digits. String must contain either escape sequence (`\\`) with digits or alphanumeric chars")
	}
	if alphabetOnly.MatchString(str) {
		return str, nil
	}

	runes := []rune(str)
	var result []rune
	var i int
	var digitCanEscape bool = false

	for i < len(runes) {
		current := runes[i]

		//Некорректный случай первой цифры
		if i == 0 && unicode.IsDigit(current) {
			return "", errors.New("string can't start with a digit")
		}

		//Обработкаэкранирования и случая \ в конце строки
		if current == '\\' {
			if i+1 < len(runes) {
				i++
				current = runes[i]
				if unicode.IsDigit(current) {
					digitCanEscape = true
				}

			} else {
				return "", errors.New("'\\' - single escape sequence can't be the last symbol in a string")
			}
		}

		if !unicode.IsDigit(current) || (unicode.IsDigit(current) && digitCanEscape) {
			//Обработка многосимвольного числа!!!
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				j := i + 1
				for j < len(runes) && unicode.IsDigit(runes[j]) {
					j++
				}
				numStr := string(runes[i+1 : j])
				number, err := strconv.Atoi(numStr)
				if err != nil {
					return "", err
				}
				if number != 1 {
					for k := 0; k < number; k++ {
						result = append(result, current)
					}
				} else {
					result = append(result, current)
				}
				i = j
			} else {
				result = append(result, current)
				i++
			}
		}

		//Ресет экранирования цифр
		if unicode.IsDigit(current) && digitCanEscape {
			digitCanEscape = false
		}
	}
	return string(result), nil
}

func generateSequence(minElements, maxElements int) string {
	const shorterChars = "abcdefghijklmnopqrstuvwxyz\\123456789"
	ancientTexts := []rune(shorterChars)
	var result []rune
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, ancientTexts[rand.Intn(len(ancientTexts)-1)])
	}
	return string(result)
}
