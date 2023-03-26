package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	fmt.Println("starting server on port 8080")

	fmt.Print("  connecting to db...")
	mySqlDb, errOpen := sql.Open("mysql", "root:my-secret-pw@(127.0.0.1:3306)/relay")
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

	db = mySqlDb

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
	"GET /users": getUsers,
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rowsUsers, errQuery := db.Query("SELECT id, email, password FROM users;")
	if errQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error querying users: %v", errQuery)
		return
	}

	var users []User
	for rowsUsers.Next() {
		var user User
		errScan := rowsUsers.Scan(&user.ID, &user.Email, &user.Password)
		if errScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error scanning user: %v", errScan)
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
