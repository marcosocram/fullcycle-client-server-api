// server.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type CotacaoServer struct {
	Bid string `json:"bid"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout para chamar a API externa
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cotacao, err := getCotacaoUSDBRL(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		log.Println("Erro ao obter cotação:", err)
		return
	}

	// Timeout para salvar no banco
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	if err := saveCotacao(ctxDB, cotacao); err != nil {
		http.Error(w, "Erro ao salvar cotação no banco", http.StatusInternalServerError)
		log.Println("Erro ao salvar cotação no banco:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}

func getCotacaoUSDBRL(ctx context.Context) (*CotacaoServer, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resultado map[string]CotacaoServer
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return nil, err
	}

	cotacao := resultado["USDBRL"]
	return &cotacao, nil
}

func saveCotacao(ctx context.Context, cotacao *CotacaoServer) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", cotacao.Bid)
	return err
}
