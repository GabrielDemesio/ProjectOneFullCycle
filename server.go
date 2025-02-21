package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func fetchDollarRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result["USDBRL"]["bid"], nil
}

func saveToDatabase(ctx context.Context, db *sql.DB, bid string) error {
	query := "INSERT INTO exchange_rates (bid, timestamp) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	_, err = stmt.ExecContext(ctx, bid, time.Now().Unix())
	return err
}

func handleQuote(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	bid, err := fetchDollarRate(ctx)
	if err != nil {
		http.Error(w, "Error fetching dollar rate: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error fetching dollar rate:", err)
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	db, err := sql.Open("sqlite3", "./exchange_rates.db")
	if err != nil {
		http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	err = saveToDatabase(dbCtx, db, bid)
	if err != nil {
		http.Error(w, "Error saving to database: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error saving to database:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"bid": bid})
	if err != nil {
		return
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./exchange_rates.db")
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS exchange_rates (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, timestamp INTEGER)")
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	http.HandleFunc("/cotacao", handleQuote)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
