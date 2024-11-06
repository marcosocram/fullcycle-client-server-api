// client.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CotacaoClient struct {
	Bid string `json:"bid"`
}

func main() {
	// Timeout para a requisição ao servidor
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cotacao, err := getCotacaoServerAPI(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := saveFileCotacao(cotacao); err != nil {
		log.Fatal("Erro ao salvar cotação no arquivo:", err)
	}

	fmt.Println("Cotação salva com sucesso!")
}

func getCotacaoServerAPI(ctx context.Context) (*CotacaoClient, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%v", string(bodyBytes))
	}

	var cotacao CotacaoClient
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func saveFileCotacao(cotacao *CotacaoClient) error {
	conteudo := fmt.Sprintf("Dólar: %s\n", cotacao.Bid)
	return ioutil.WriteFile("cotacao.txt", []byte(conteudo), 0644)
}
