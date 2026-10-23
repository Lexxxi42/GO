package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lab01/src"
	"log"
	"net/http"
	"os"
)

const apiURL = "http://127.0.0.1:1234/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func readToken() string {
	return os.Getenv("HF_TOKEN")
}

func readDataset(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Ошибка открытия файла: ", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Ошибка чтения файла: ", err)
	}
	return string(data)
}

func buildRequest(model string, messages []Message) []byte {
	request := RequestBody{
		Model:       model,
		Messages:    messages,
		MaxTokens:   500,
		Temperature: 0.5,
	}
	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func sendRequest(requestBody []byte) []byte {
	request, err := http.NewRequest(
		"POST",
		apiURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("API error: %s\n%s", response.Status, string(body))
	}
	return body
}

func parseResponse(body []byte) string {
	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content
	}
	return ""
}

func chooseModel() string {
	fmt.Println("Выберите модель:")
	fmt.Println("1. bonsai-8b")
	fmt.Println("2. qwen3-1.7b")
	fmt.Println("3. Qwen3-8B")

	var choice int
	fmt.Scan(&choice)

	models := []string{
		"bonsai-8b",
		"qwen/qwen3-1.7b",
		"ollama/qwen3",
	}
	if choice < 1 || choice > len(models) {
		log.Fatal("Нет такой модели")
	}
	return models[choice-1]
}

func ask(model string, context []Message, question string) []Message {
	context = append(context, Message{Role: "user", Content: question})

	response := sendRequest(buildRequest(model, context))
	answer := parseResponse(response)

	fmt.Println("Ответ:", answer)
	context = append(context, Message{Role: "assistant", Content: answer})

	return context
}

func main() {
	model := chooseModel()

	system := Message{
		Role:    "system",
		Content: "Ты — эксперт-разработчик в GO.",
	}
	context := []Message{system}

	context = ask(model, context, "Привет, что такое go?")

	health := src.ReadCSV[src.Insurance]("data/Health_Insurance.csv")
	avgAge := src.CountByRegion(health)
	var regionsText string
	for region, count := range avgAge {
		regionsText += fmt.Sprintf("%s: %d, ", region, count)
	}
	dataMessage := Message{
		Role:    "system",
		Content: fmt.Sprintf("Вот актуальная статистика. Всего страховок: %d. По регионам: %s", regionsText),
	}
	context = append(context, dataMessage)
	context = ask(model, context, "Сколько всего страховок?")
}
