package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/desafio/sever/api/model"
	"github.com/desafio/sever/db"
)

var (
	API_URL           = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	METHOD            = "GET"
	API_TIMEOUT int32 = 180
	CTX_TIMEOUT int32 = 2
)

func Handler(w http.ResponseWriter, r *http.Request) {
	response, err := callApi()

	if err != nil {
		ResponseJSON(w, "", errors.New("timeout"))
		return
	}

	if err != nil {
		ResponseJSON(w, "", errors.New("the api took more than 200ms to respond"))
		return
	}

	conn, err := db.ConnectSqlite3()
	if err = conn.CreateTable(); err != nil {
		log.Println("could not create table sqlite3 error message ->", err)
	}
	go func() {
		conn.Save(response)
	}()

	ResponseJSON(w, response.Bid, nil)
}

func callApi() (model.Awesomeapi, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CTX_TIMEOUT)*time.Second)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, METHOD, API_URL, nil)
	if err != nil {
		log.Println(err)
		return model.Awesomeapi{}, errors.New("fail on mount request with context")
	}

	req, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println(err)
		return model.Awesomeapi{}, errors.New("error to execute request")
	}
	defer req.Body.Close()

	bodyResponse, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return model.Awesomeapi{}, errors.New("error to read request body")
	}

	var response model.USDBRL
	if err := json.Unmarshal(bodyResponse, &response); err != nil {
		log.Println(err)
		return model.Awesomeapi{}, errors.New("error to serialize request response")
	}

	elapsed := time.Since(start).Milliseconds()
	defer fmt.Printf("callApi took %d (ms) to execute\n", elapsed)

	select {
	case <-ctx.Done():
		return model.Awesomeapi{}, errors.New("timeout error to call external api")
	case <-time.After(time.Duration(API_TIMEOUT) * time.Millisecond):
		return response.USDBRL, nil
	}
}

// same func without ctx timeout
// func callApi() (Awesomeapi, error) {
// 	start := time.Now()

// 	r, err := http.NewRequest(METHOD, API_URL, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return Awesomeapi{}, errors.New("falha ao montar request com context")
// 	}

// 	req, err := http.DefaultClient.Do(r)
// 	if err != nil {
// 		log.Println(err)
// 		return Awesomeapi{}, errors.New("erro ao executar request")
// 	}
// 	defer req.Body.Close()

// 	bodyResponse, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		log.Println(err)
// 		return Awesomeapi{}, errors.New("erro ao ler corpo da requisicao")
// 	}

// 	var response USDBRL
// 	if err := json.Unmarshal(bodyResponse, &response); err != nil {
// 		log.Println(err)
// 		return Awesomeapi{}, errors.New("erro ao na serializacao dos dados retornados")
// 	}

// 	elapsed := time.Since(start).Milliseconds()
// 	fmt.Println(elapsed)

// 	return response.USDBRL, nil
// }
