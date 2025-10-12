/*
Реализовать упрощённый аналог UNIX-утилиты sort (сортировка строк).

Программа должна читать строки (из файла или STDIN) и выводить их отсортированными.

Обязательные флаги (как в GNU sort):
-k N — сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию).
Например, «sort -k 2» отсортирует строки по второму столбцу каждой строки.
-n — сортировать по числовому значению (строки интерпретируются как числа).
-r — сортировать в обратном порядке (reverse).
-u — не выводить повторяющиеся строки (только уникальные).

Дополнительные флаги:
-M — сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат.
-b — игнорировать хвостовые пробелы (trailing blanks).
-c — проверить, отсортированы ли данные; если нет, вывести сообщение об этом.
-h — сортировать по числовому значению с учётом суффиксов (например, К = килобайт, М = мегабайт — человекочитаемые размеры).

Программа должна корректно обрабатывать комбинации флагов (например, -nr — числовая сортировка в обратном порядке, и т.д.).
Необходимо предусмотреть эффективную обработку больших файлов.
Код должен проходить все тесты, а также проверки go vet и golint
(понимание, что требуются надлежащие комментарии, имена и структура программы).
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type sortConfig struct {
	filePath               string
	column                 int
	delimiter              string
	isTable                bool
	isNumeric              bool //sort
	isReverse              bool //postsort
	isUnique               bool //presort
	isMonth                bool //sort
	isIgnoreTrailingBlanks bool //presort
	isSortedCheck          bool //postsort
	isSuffixEnabled        bool //presort
}

func parseConfig() sortConfig {
	filePath := flag.String("file", "", "Path to file with data")
	column := flag.Int("k", 0, "Column number to sort by")
	delimiter := flag.String("N", "\t", "Column delimiter")
	isTable := flag.Bool("T", true, "Incomming data is in table format")
	isNumeric := flag.Bool("n", false, "Sort by numeric value")
	isReverse := flag.Bool("r", false, "Sort descending (default - ascending)")
	isUnique := flag.Bool("u", false, "Show unique rows")
	isMonth := flag.Bool("M", false, "Sort by month name")
	isIgnoreTrailingBlanks := flag.Bool("b", false, "Ignore trailing blanks")
	isSortedCheck := flag.Bool("c", false, "Check if output is sorted")
	isSuffixEnabled := flag.Bool("h", false, "Check values based on suffix")
	flag.Parse()

	return sortConfig{
		filePath:               *filePath,
		column:                 *column,
		delimiter:              *delimiter,
		isTable:                *isTable,
		isNumeric:              *isNumeric,
		isReverse:              *isReverse,
		isUnique:               *isUnique,
		isMonth:                *isMonth,
		isIgnoreTrailingBlanks: *isIgnoreTrailingBlanks,
		isSortedCheck:          *isSortedCheck,
		isSuffixEnabled:        *isSuffixEnabled,
	}
}

func main() {
	config := parseConfig()
	if config.isMonth && config.isNumeric {
		RaiseErrorAndStop(errors.New("month and Numeric flags can't be used at the same time"))
	}
	if config.isSuffixEnabled && config.isNumeric {
		RaiseErrorAndStop(errors.New("suffix and Numeric flags can't be used at the same time"))
	}
	if !config.isTable && config.column != 0 {
		RaiseErrorAndStop(errors.New("non-Table data and column flags can't be used at the same time"))
	}

	var data []string
	var err error

	if config.filePath != "" {
		data, err = ReadFile(config.filePath)
		if err != nil {
			RaiseErrorAndStop(err)
		}
	} else {
		data = GenerateSequence(10, 20)
	}
	var dataInTableForm [][]string
	for _, v := range data {
		dataInTableForm = append(dataInTableForm, strings.Split(v, config.delimiter))
	}
	ShowDataTable(dataInTableForm)

	dataInTableForm = PreSortActions(config, dataInTableForm)

	//SORTING

	dataInTableForm = PostSortActions(config, dataInTableForm)
	ShowDataTable(dataInTableForm)
}

func ShowDataTable(data [][]string) {
	for _, line := range data {
		for _, val := range line {
			fmt.Printf("%s ", val)
		}
		fmt.Print("\n")
	}
}

func RaiseErrorAndStop(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func SortTable(config sortConfig, data [][]string) [][]string {
	data = PreSortActions(config, data)

	//Main sorting

	data = PostSortActions(config, data)
	return data
}

func PostSortActions(config sortConfig, data [][]string) [][]string {
	if config.isReverse {
		data = reverse(data)
	}
	if config.isSortedCheck {
		if !isSorted(data, config.column) {
			fmt.Println("Data was not sorted properly:", config.column, data[0][config.column])
		}
	}
	return data
}

func reverse(data [][]string) [][]string {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func isSorted(data [][]string, column int) bool {
	for i := 1; i < len(data); i++ {
		if data[i-1][column] > data[i][column] {
			return false
		}
	}
	return true
}

func PreSortActions(config sortConfig, data [][]string) [][]string {
	if config.isIgnoreTrailingBlanks {
		data = trimTrailingBlanks(data)
	}
	if config.isUnique {
		data = unique(data, config.column)
	}
	if config.isSuffixEnabled {
		data = suffixConversion(data, config.column)
	}
	return data
}

func suffixConversion(data [][]string, column int) [][]string {
	// TODO: Реализовать считывание суффиксов (Преобразовать число с суффиксом в число или сравнивать суффиксы, а потом числа)
	return data
}

func trimTrailingBlanks(data [][]string) [][]string {
	result := make([][]string, 0)
	for _, lines := range data {
		line := make([]string, 0)
		for _, val := range lines {
			line = append(lines, strings.TrimRightFunc(val, unicode.IsSpace))
		}
		result = append(result, line)
	}
	return result
}

func unique(data [][]string, column int) [][]string {
	exist := make(map[string]bool)
	unique := make([][]string, 0)
	for _, lineval := range data {
		key := fmt.Sprintf("%v", lineval)
		if !exist[key] {
			exist[key] = true
			unique = append(unique, lineval)

		}
	}
	return unique
}

func ReadFile(filePath string) ([]string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(file), "\n"), nil
}

func GenerateSequence(minElements, maxElements int) []string {
	var result []string
	/*
		TODO: Сгенерить инфу: id, datetime, month, word, word with trailing blanks, suffix, float64 value
	*/
	return result
}
