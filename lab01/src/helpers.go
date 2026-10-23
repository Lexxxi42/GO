package src

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

// yearOf парсит строку года в int, при ошибке возвращает 0.
func yearOf(s string) int {
	y, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return y
}

// OpenFile открывает файл и паникует при ошибке.
func OpenFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic("Ошибка открытия файла: " + path)
	}
	return f
}

// ReadJSON открывает JSON и возвращает коллекцию элементов.
func ReadJSON[T any](path string) T {
	file := OpenFile(path)
	defer file.Close()
	var data T
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		panic("Ошибка парсинга JSON: " + err.Error())
	}
	return data
}

// ReadCSV считывает игры из CSV и возвращает срез Vgsale.
func ReadCSV[T any](path string) []T {
	file := OpenFile(path)
	defer file.Close()
	var data []T
	if err := gocsv.UnmarshalFile(file, &data); err != nil {
		panic("Ошибка парсинга CSV: " + err.Error())
	}
	fmt.Printf("Прочитано: %d строк\n", len(data))
	return data
}

// RowDesc описывает, как печатать любой список: заголовки, ширины колонок и
// функция записи -> ячейки. Table=false печатает строки подряд.
type RowDesc[T any] struct {
	Title   string
	Headers []string
	Widths  []int
	Cell    func(T) []string
	Table   bool
}

// PrintRows печатает data согласно описанию d. Используется для всех датасетов.
func PrintRows[T any](data []T, d RowDesc[T]) {
	if d.Title != "" {
		fmt.Println(d.Title)
	}
	if d.Table && len(d.Headers) > 0 {
		var b strings.Builder
		for i, h := range d.Headers {
			b.WriteString(fmt.Sprintf("|%-*s", d.Widths[i], h))
		}
		fmt.Println(b.String() + "|")
		fmt.Println(strings.Repeat("-", b.Len()+1))
		for _, v := range data {
			cs := d.Cell(v)
			b.Reset()
			for i, c := range cs {
				b.WriteString(fmt.Sprintf("|%-*s", d.Widths[i], c))
			}
			fmt.Println(b.String() + "|")
		}
		return
	}
	for _, v := range data {
		fmt.Println(strings.Join(d.Cell(v), "  "))
	}
}

// table собирает описание таблицы, чтобы не повторять структуру RowDesc.
func table[T any](title string, headers []string, widths []int, cell func(T) []string) RowDesc[T] {
	return RowDesc[T]{Title: title, Headers: headers, Widths: widths, Cell: cell, Table: true}
}

// Filter возвращает элементы data, для которых keep() == true.
func Filter[T any](data []T, keep func(T) bool) []T {
	var out []T
	for _, v := range data {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// CountBy группирует data по ключу и возвращает счётчики.
func CountBy[T any, K comparable](data []T, key func(T) K) map[K]int {
	m := make(map[K]int)
	for _, v := range data {
		m[key(v)]++
	}
	return m
}

// TopBy сортирует data по less (первый = "лучший") и возвращает первые n.
func TopBy[T any](data []T, n int, less func(a, b T) bool) []T {
	out := append([]T(nil), data...)
	sort.SliceStable(out, func(i, j int) bool { return less(out[i], out[j]) })
	if n > len(out) {
		n = len(out)
	}
	return out[:n]
}

func Chapter(n int) {
	var m = 10
	var delimeter = strings.Repeat("=", m*2)
	var tab = strings.Repeat(" ", m/2)
	fmt.Printf("%s\n%sChapter %d\n%s\n", delimeter, tab, n, delimeter)
}
