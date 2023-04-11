package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

var DB *sql.DB

var OPEN_AI_SK = "sk-"

func PostApi(w http.ResponseWriter, r *http.Request) {
	var reqBody map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&reqBody)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", errDecode)
		return
	}

	reqCost := float64(len(reqBody["prompt"].(string))) / 4
	fmt.Printf("  prompt tokens: %.2f\n", reqCost)

	targetBody, errEncode := json.Marshal(reqBody)
	if errEncode != nil {
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		return
	}

	targetReq, err := http.NewRequest(r.Method, "https://api.openai.com"+r.RequestURI[11:], bytes.NewBuffer(targetBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetReq.Header.Set("Authorization", "Bearer "+OPEN_AI_SK)
	targetReq.Header.Set("Content-Type", "application/json")
	targetClient := http.Client{}

	targetResp, err := targetClient.Do(targetReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	flusher := w.(http.Flusher)

	for {
		data := make([]byte, 1024)
		_, err := targetResp.Body.Read(data)
		if err != nil {
			fmt.Printf("  prompt + completion tokens: %.2f\n", reqCost)
			email, _, _ := r.BasicAuth()
			_, errExec := DB.Exec("UPDATE users SET number_of_tokens = number_of_tokens - ? WHERE email = ?;", reqCost, email)
			if errExec != nil {
				fmt.Fprintf(w, "  error decrementing user tokens: %v", errExec)
			}
			return
		}

		reqCost += 1
		w.Write(data)
		flusher.Flush()
	}
}
