package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI)

	targetReq, err := http.NewRequest(r.Method, "https://api.openai.com"+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetReq.Header = r.Header
	targetClient := http.Client{}

	targetResp, err := targetClient.Do(targetReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer targetResp.Body.Close()

	targetData, err := ioutil.ReadAll(targetResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the response to the original client
	for key, values := range targetResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(targetResp.StatusCode)
	w.Write(targetData)
}
