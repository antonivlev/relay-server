package api

import (
	"net/http"
)

var OPEN_AI_SK = "sk-..."

func PostApi(w http.ResponseWriter, r *http.Request) {
	targetReq, err := http.NewRequest(r.Method, "https://api.openai.com"+r.RequestURI[11:], r.Body)
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
			return
		}

		w.Write(data)
		flusher.Flush()
	}
}
