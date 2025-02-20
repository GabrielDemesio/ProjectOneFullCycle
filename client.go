package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func fetchCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result["bid"], nil
}

func saveToFile(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)
	return ioutil.WriteFile("cotacao.txt", []byte(content), 0644)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	bid, err := fetchCotacao(ctx)
	if err != nil {
		log.Println("Error fetching cotação:", err)
		return
	}

	err = saveToFile(bid)
	if err != nil {
		log.Println("Error saving to file:", err)
		return
	}

	fmt.Println("Cotação salva em 'cotacao.txt'")
}
