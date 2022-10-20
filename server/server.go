package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"time"
)

var api_url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

var db_path = "./cotacao.db"

type Quote struct {
	Bid string `json:"bid"`
}

func main() {
	err := createTable()
	if err != nil {
		panic(err)
	}
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Call /cotacao to receive the current USD-BRL exchange rate"))
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.Get(api_url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("Request timed out"))
	default:
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		m := make(map[string]Quote)
		err = json.Unmarshal(body, &m)
		if err != nil {
			w.Write([]byte("Internal server error"))
		}

		var cs []Quote
		for _, v := range m {
			cs = append(cs, v)
		}

		err = persist(cs[0].Bid, r.RemoteAddr)
		if err != nil {
			w.Write([]byte("Internal server error"))
		}

		out := `{"bid": "` + cs[0].Bid + `"}`
		w.Write([]byte(out))
	}
}

func persist(bid, remotePath string) error {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	sqlStmt := `insert into cotacao(bid, remote) values(?, ?);`

	_, err = db.ExecContext(ctx, sqlStmt, bid, remotePath)
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}

func createTable() error {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `create table if not exists cotacao (id integer not null primary key, bid text, remote text);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
