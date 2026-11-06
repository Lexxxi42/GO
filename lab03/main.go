package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
)

// fileNameRe — разрешённое имя файла. По контракту: латиница/цифры/_/- и строго .html.
// Никаких путей, слэшей и точек — иначе можно вычитать что попало за пределами набора.
var fileNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+\.html$`)

// linkRe ищет открывающий тег <a> с внешним href (http:// или https://).
// Закрывающий </a> не проверяем: контракт велит считать по открывающему тегу.
// Здесь хитро, но работает — регулярка жадная до href, остальное ей неважно.
var linkRe = regexp.MustCompile(`<a\s+[^>]*href="https?://[^"]*"`)

// writeJSON сериализует body в JSON и пишет в response.
// input: w — ответ, code — HTTP-код, body — данные (map или struct)
// output: JSON в теле ответа
func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println("Ошибка записи JSON:", err)
	}
}

// countExternalLinks считает внешние ссылки в HTML-файле.
// input: filename — имя файла из предопределённого набора
// output: количество внешних ссылок (int)
func countExternalLinks(filename string) int {
	file, err := os.ReadFile(filename)
	if err != nil {
		// сюда не дойдём: хендлер заранее проверил файл и вернул 404
		log.Fatal("Не удалось прочитать файл: " + filename)
	}
	return len(linkRe.FindAllString(string(file), -1))
}

// getAHrefsCounts — GET /a_count?filename=page1.html
// Возможные ответы: 200 (количество ссылок), 400 (плохое имя), 404 (файла нет).
func getAHrefsCounts(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	// невалидное имя — сразу отбой, даже не смотрим на диск
	if !fileNameRe.MatchString(filename) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_filename",
			"message": "Некорректный запрос. Неверный формат имени файла или использование запрещенных символов (пути, точки).",
		})
		return
	}

	// файл есть? если нет — 404 по контракту
	if _, err := os.Stat(filename); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "file_not_found",
			"message": "Файл с указанным именем не найден в предопределенном наборе.",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"filename":             filename,
		"external_links_count": countExternalLinks(filename),
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/a_count", getAHrefsCounts).Methods(http.MethodGet)
	r.HandleFunc("/phones", getPhonesCount).Methods(http.MethodGet)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}