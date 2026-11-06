package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
)

// txtNameRe — то же ограничение, что и для html, но расширение должно быть .txt.
var txtNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+\.txt$`)

// phoneRe ищет российский телефон: 8 или +7, дальше 10 цифр, между которыми
// могут стоять пробелы, скобки и дефисы: 8 (912) 345-67-89, +7 922 111-22-33.
// Границы (чтобы не схватить кусок длинного числа) проверяем отдельно — см.
// isValidPhoneMatch, потому что в Go regexp нет lookahead/lookbehind.
var phoneRe = regexp.MustCompile(`(?:8|\+7)[\s(-]*\d{3}[\s)-]*\d{3}[\s-]*\d{2}[\s-]*\d{2}`)

// isDigitByte — это цифра?
func isDigitByte(c byte) bool {
	return c >= '0' && c <= '9'
}

// isValidPhoneMatch проверяет, что вокруг номера нет других цифр. Иначе это
// кусок длинного числа вроде 123456789012345, а не телефон.
// input: text — весь файл, start/end — границы найденного совпадения
// output: true, если по соседству нет цифр
func isValidPhoneMatch(text string, start, end int) bool {
	if start > 0 && isDigitByte(text[start-1]) {
		return false
	}
	if end < len(text) && isDigitByte(text[end]) {
		return false
	}
	return true
}

// getPhonesCount — GET /phones?filename=phones.txt
// Работает как /a_count, но по .txt и считает телефоны.
// Ответ: filename, phones_count и first_position (позиция первого номера, -1 если нет).
func getPhonesCount(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	if !txtNameRe.MatchString(filename) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_filename",
			"message": "Некорректный запрос. Неверный формат имени файла или использование запрещенных символов (пути, точки).",
		})
		return
	}

	if _, err := os.Stat(filename); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "file_not_found",
			"message": "Файл с указанным именем не найден в предопределенном наборе.",
		})
		return
	}

	// читаем файл один раз, позиции считаем по этой строке
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Не удалось прочитать файл: " + filename)
	}
	text := string(file)

	// ищем все совпадения, но в счёт идут только с "чистыми" границами
	matches := phoneRe.FindAllStringIndex(text, -1)
	count := 0
	firstPosition := -1
	for _, m := range matches {
		if !isValidPhoneMatch(text, m[0], m[1]) {
			continue
		}
		if firstPosition == -1 {
			firstPosition = m[0]
		}
		count++
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"filename":       filename,
		"phones_count":   count,
		"first_position": firstPosition,
	})
}