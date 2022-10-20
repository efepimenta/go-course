package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

var server_url = "http://127.0.0.1:8080/cotacao"

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	c, err := callApi()
	if err != nil {
		panic(err)
	}
	writeFile(c)
}

func callApi() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", server_url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func writeFile(c *Cotacao) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(`Dólar: {` + c.Bid + `}`))
	if err != nil {
		panic(err)
	}
	return nil
}
