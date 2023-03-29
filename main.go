package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/antonivlev/relay-server/api"
	"github.com/antonivlev/relay-server/users"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("starting server on port 8080")

	fmt.Print("  connecting to db...")
	mySqlDb, errOpen := sql.Open("mysql", "root:my-secret-pw@(127.0.0.1:3306)/relay?parseTime=true")
	if errOpen != nil {
		log.Fatal(errOpen.Error())
	}
	fmt.Print("done\n")

	fmt.Print("  pinging db...")
	errPing := mySqlDb.Ping()
	if errPing != nil {
		log.Fatal(errPing.Error())
	}
	fmt.Print("done\n")

	users.DB = mySqlDb

	fmt.Println("server running")
	errServe := http.ListenAndServe(":8080", http.HandlerFunc(handlerOfAllRequests))
	log.Fatal(errServe)
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

var routeToHandler = map[string]func(http.ResponseWriter, *http.Request){
	"POST /login": users.PostLogin,
	"GET /users":  users.GetUsers,
	"POST /users": users.PostUsers,
	"POST /api/*": api.PostApi,
}
