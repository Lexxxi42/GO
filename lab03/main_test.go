package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// doRequest дергает хендлер с query-параметром и возвращает записанный ответ.
// input: query — строка вида "filename=page1.html"
// output: recorder с кодом, заголовками и телом
func doRequest(query string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/a_count?"+query, nil)
	rec := httptest.NewRecorder()
	getAHrefsCounts(rec, req)
	return rec
}

// TestGetAHrefsCounts_OK — стандартный успешный ответ от тестировщика.
// Файл из задания page1.html: одна внешняя ссылка https://giga.chat.
func TestGetAHrefsCounts_OK(t *testing.T) {
	rec := doRequest("filename=page1.html")

	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, ожидался application/json", ct)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("тело не JSON: %v; тело: %s", err, rec.Body.String())
	}
	if resp["filename"] != "page1.html" {
		t.Errorf("filename = %v, ожидался page1.html", resp["filename"])
	}
	if resp["external_links_count"] != float64(1) {
		t.Errorf("external_links_count = %v, ожидалось 1", resp["external_links_count"])
	}
}

// TestGetAHrefsCounts_400 — невалидное имя: точки, слэши, null.
func TestGetAHrefsCounts_400(t *testing.T) {
	for _, bad := range []string{"../x", "/etc/passwd", "page1", "page1.htm", "", "file.html?a=1"} {
		rec := doRequest("filename=" + bad)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("filename=%q: код = %d, ожидался 400", bad, rec.Code)
		}
	}
}

// TestGetAHrefsCounts_404 — имя валидное по паттерну, но файла нет.
func TestGetAHrefsCounts_404(t *testing.T) {
	rec := doRequest("filename=page9.html")
	if rec.Code != http.StatusNotFound {
		t.Errorf("код = %d, ожидался 404", rec.Code)
	}
}

// TestGetAHrefsCounts_NoLinks — page2.html: ссылок нет вообще.
// Файл создаём на лету, чтобы тест не зависел от чужих файлов.
func TestGetAHrefsCounts_NoLinks(t *testing.T) {
	if err := os.WriteFile("page2.html", []byte("<p>тут нет ссылок</p>"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("page2.html")

	rec := doRequest("filename=page2.html")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"external_links_count":0`) {
		t.Errorf("ожидался 0, получили: %s", rec.Body.String())
	}
}

// TestGetAHrefsCounts_Tricky — page3.html: относительные, mailto, tel и https3 не считаются.
func TestGetAHrefsCounts_Tricky(t *testing.T) {
	rec := doRequest("filename=page3.html")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	// в page3.html внешних ссылок ровно две: https://ya.ru и http://habr.com
	if !strings.Contains(rec.Body.String(), `"external_links_count":2`) {
		t.Errorf("ожидался 2, получили: %s", rec.Body.String())
	}
}

// TestGetAHrefsCounts_UnclosedTags — незакрытый <a> с внешним href должен учитываться.
func TestGetAHrefsCounts_UnclosedTags(t *testing.T) {
	if err := os.WriteFile("page_unclosed.html", []byte(
		`<a href="https://example.com">без закрывающего тега`), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("page_unclosed.html")

	rec := doRequest("filename=page_unclosed.html")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"external_links_count":1`) {
		t.Errorf("ожидался 1, получили: %s", rec.Body.String())
	}
}