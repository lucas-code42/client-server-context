package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	RUN_TIMES      = 100
	DELAY          = 80
	SERVER_TIMEOUT = 2
	CTX_TIMEOUT    = 100
	SERVER_METHOD  = "GET"
	SERVER_HOST    = "http://localhost:8080/cotacao"
)

type ServerResponse struct {
	Dolar string `json:"dolar,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

func callServer() (ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(SERVER_TIMEOUT)*time.Second)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, SERVER_METHOD, SERVER_HOST, nil)
	if err != nil {
		log.Println("error to mount request")
		log.Fatal(err)
	}

	req, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println("error to execute request")
		log.Fatal(err)
	}
	defer req.Body.Close()

	bodyResponse, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("error to read request response")
		log.Fatal(err)
	}

	var response ServerResponse
	if err = json.Unmarshal(bodyResponse, &response); err != nil {
		log.Println("error to serialize response")
		log.Fatal(err)
	}

	select {
	case <-ctx.Done():
		return ServerResponse{}, errors.New("server timeout")
	case <-time.After(time.Duration(CTX_TIMEOUT) * time.Millisecond):
		return response, nil
	}
}

func writeFile(data string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return errors.New("could not create file -> " + err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(data + "\n")
	if err != nil {
		return errors.New("could write on file -> " + err.Error())
	}

	return nil
}

func main() {
	fmt.Println("client is running...")

	for {
		RUN_TIMES -= 1
		time.Sleep(time.Duration(DELAY) * time.Millisecond)

		start := time.Now()
		r, err := callServer()
		if err != nil {
			log.Fatal("could not call server ->", err)
		}
		elapsed := time.Since(start).Milliseconds()

		if r.Dolar != "" {
			if err = writeFile("Dolar: " + r.Dolar); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("server took %d ms to respond\n", elapsed)
		fmt.Println("err:", err)
		fmt.Println("response:", r)

		if RUN_TIMES == 0 {
			break
		}
	}

}
