package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	valor, err := valorDolar()
	if err != nil {
		panic(err)
	}

	escreveCotacao(valor)
}

func valorDolar() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func escreveCotacao(valor string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(fmt.Sprintf("Dólar: %s\n", valor))
	if err != nil {
		panic(err)
	}

	f.Close()
}
