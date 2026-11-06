package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// doPhonesRequest дергает хендлер /phones и возвращает записанный ответ.
func doPhonesRequest(query string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/phones?"+query, nil)
	rec := httptest.NewRecorder()
	getPhonesCount(rec, req)
	return rec
}

// TestGetPhonesCount_OK — в phones.txt три номера (8, +7 и 8-800), первый на позиции 13.
func TestGetPhonesCount_OK(t *testing.T) {
	rec := doPhonesRequest("filename=phones.txt")

	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("тело не JSON: %v; тело: %s", err, rec.Body.String())
	}
	if resp["filename"] != "phones.txt" {
		t.Errorf("filename = %v, ожидался phones.txt", resp["filename"])
	}
	if resp["phones_count"] != float64(3) {
		t.Errorf("phones_count = %v, ожидалось 3", resp["phones_count"])
	}
	// "Звони по номеру " — 16 символов: 13 кириллических (26 байт в UTF-8) + 3 пробела = 29.
	if resp["first_position"] != float64(29) {
		t.Errorf("first_position = %v, ожидалась 29", resp["first_position"])
	}
}

// TestGetPhonesCount_NoPhones — файл без номеров: count=0, first_position=-1.
func TestGetPhonesCount_NoPhones(t *testing.T) {
	if err := os.WriteFile("empty.txt", []byte("тут одни буквы, даже цифр нет\n"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("empty.txt")

	rec := doPhonesRequest("filename=empty.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"phones_count":0`) ||
		!strings.Contains(rec.Body.String(), `"first_position":-1`) {
		t.Errorf("ожидался 0 и -1, получили: %s", rec.Body.String())
	}
}

// TestGetPhonesCount_400 — плохие имена файлов: пути, точки, не .txt.
func TestGetPhonesCount_400(t *testing.T) {
	for _, bad := range []string{"../secret.txt", "phones", "phones.html", "a/b.txt", "..\\x.txt", ""} {
		rec := doPhonesRequest("filename=" + bad)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("filename=%q: код = %d, ожидался 400", bad, rec.Code)
		}
	}
}

// TestGetPhonesCount_404 — имя валидное, файла нет.
func TestGetPhonesCount_404(t *testing.T) {
	rec := doPhonesRequest("filename=missing.txt")
	if rec.Code != http.StatusNotFound {
		t.Errorf("код = %d, ожидался 404", rec.Code)
	}
}

// TestGetPhonesCount_NotANumber — длинное число не считается телефоном:
// 123456789012345 под маску не подходит.
func TestGetPhonesCount_NotANumber(t *testing.T) {
	if err := os.WriteFile("big.txt", []byte("идентификатор 123456789012345 тут"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("big.txt")

	rec := doPhonesRequest("filename=big.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"phones_count":0`) {
		t.Errorf("длинное число посчиталось номером: %s", rec.Body.String())
	}
}

// TestGetPhonesCount_CyclicPrefix — номера, идущие подряд через пробел, считаются оба.
func TestGetPhonesCount_CyclicPrefix(t *testing.T) {
	if err := os.WriteFile("two.txt", []byte("8 555 111-22-33 8 555 444-55-66"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("two.txt")

	rec := doPhonesRequest("filename=two.txt")
	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"phones_count":2`) {
		t.Errorf("ожидалось 2 номера, получили: %s", rec.Body.String())
	}
}