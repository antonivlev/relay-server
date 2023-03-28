package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

var DB *sql.DB

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rowsUsers, errQuery := DB.Query("SELECT id, email, password FROM users;")
	if errQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errQuery)
		return
	}

	var users []User
	for rowsUsers.Next() {
		var user User
		errScan := rowsUsers.Scan(&user.ID, &user.Email, &user.Password)
		if errScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", errScan)
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func PostUsers(w http.ResponseWriter, r *http.Request) {
	var user User
	errDecode := json.NewDecoder(r.Body).Decode(&user)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", errDecode)
		return
	}

	result, errExec := DB.Exec("INSERT INTO users (created_at, email, password) VALUES (now(), ?, ?);", user.Email, user.Password)
	if errExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errExec)
		return
	}

	newUserId, errLastInsertId := result.LastInsertId()
	if errLastInsertId != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errLastInsertId)
		return
	}

	var newUser User
	userRow := DB.QueryRow("SELECT id, created_at, email, password FROM users WHERE id = ?;", newUserId)
	errScan := userRow.Scan(&newUser.ID, &newUser.CreatedAt, &newUser.Email, &newUser.Password)
	if errScan != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errScan)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO
}
