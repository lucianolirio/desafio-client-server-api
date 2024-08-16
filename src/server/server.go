package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type Moeda struct {
	USDBRL Cotacao `json:"USDBRL"`
}

type Cotacao struct {
	ID         uint `gorm:"primaryKey"`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

func main() {
	http.HandleFunc("/cotacao", cotacaoDolarHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoDolarHandler(w http.ResponseWriter, r *http.Request) {
	moeda, err := cotacaoDolar()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
	} else {
		err = salvarCotacao(&moeda.USDBRL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(moeda.USDBRL.Bid))
		}
	}
}

func cotacaoDolar() (*Moeda, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cambio Moeda
	err = json.Unmarshal(body, &cambio)
	if err != nil {
		return nil, err
	}
	return &cambio, nil
}

func salvarCotacao(cocacao *Cotacao) error {

	db, err := gorm.Open(sqlite.Dialector{
		DSN:        "cotacao.db",
		DriverName: "sqlite",
	}, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Cotacao{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err = db.WithContext(ctx).Create(&cocacao).Error
	if err != nil {
		panic(err)
	}
	return nil
}
