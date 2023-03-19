package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting server on port 8080")

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handlerOfAllRequests)))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, these are users"))
}

var routeToHandler = map[string]func(http.ResponseWriter, *http.Request){
	"GET /users": getUsers,
}

func handlerOfAllRequests(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI)

	routeHandler := routeToHandler[r.Method+" "+r.RequestURI]
	if routeHandler != nil {
		routeHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
