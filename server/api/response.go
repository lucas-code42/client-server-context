package api

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, data string, err error) {
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		body := map[string]string{"msg": err.Error()}
		w.Write(parseJSON(body))
		return
	}

	body := map[string]string{"dolar": data}
	w.Write(parseJSON(body))
}

func parseJSON(body map[string]string) []byte {
	JSONBody, err := json.Marshal(body)
	if err != nil {
		return []byte("server error")
	}
	return JSONBody
}
