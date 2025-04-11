package web

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// формирует URL передавая те же GET параметры, что были переданы в запросе
func urlWithParameters(r *http.Request, urlStr string) (string, error) {

	// Формируем URL для запроса к PostgREST с параметрами
	postgrestURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Извлекаем параметры GET-запроса
	params := r.URL.Query()
	// Добавляем параметры к URL
	postgrestURL.RawQuery = params.Encode()

	return postgrestURL.String(), nil
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]

	postgrestURL, err := urlWithParameters(r, "http://localhost:3000/"+resource)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	// Выполняем GET-запрос к PostgREST
	req, err := http.NewRequest("GET", postgrestURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из исходного запроса
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to PostgREST", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проксируем ответ от PostgREST обратно клиенту
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]

	// Копируем тело запроса
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Формируем URL для запроса к PostgREST
	postgrestURL := "http://localhost:3000/" + resource

	// Выполняем POST-запрос к PostgREST
	req, err := http.NewRequest("POST", postgrestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из исходного запроса
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to PostgREST", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проксируем ответ от PostgREST обратно клиенту
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]

	// Извлекаем параметры GET-запроса (если есть)
	params := r.URL.Query()

	// Формируем URL для запроса к PostgREST с параметрами
	postgrestURL, err := url.Parse("http://localhost:3000/" + resource)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	// Добавляем параметры к URL
	postgrestURL.RawQuery = params.Encode()

	// Выполняем DELETE-запрос к PostgREST
	req, err := http.NewRequest("DELETE", postgrestURL.String(), nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из исходного запроса
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to PostgREST", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проксируем ответ от PostgREST обратно клиенту
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]

	// Копируем тело запроса
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	postgrestURL, err := urlWithParameters(r, "http://localhost:3000/"+resource)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	// Выполняем PATCH-запрос к PostgREST
	req, err := http.NewRequest("PATCH", postgrestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из исходного запроса
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to PostgREST", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Проксируем ответ от PostgREST обратно клиенту
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
