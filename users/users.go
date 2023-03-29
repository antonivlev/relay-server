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
	Password  string    `json:"-"`
}

var DB *sql.DB

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rowsUsers, errQuery := DB.Query("SELECT id, created_at, email FROM users;")
	if errQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errQuery)
		return
	}

	var users []User
	for rowsUsers.Next() {
		var user User
		errScan := rowsUsers.Scan(&user.ID, &user.CreatedAt, &user.Email)
		if errScan != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", errScan)
			return
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
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
	userRow := DB.QueryRow("SELECT id, created_at, email FROM users WHERE id = ?;", newUserId)
	errScan := userRow.Scan(&newUser.ID, &newUser.CreatedAt, &newUser.Email)
	if errScan != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errScan)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var userTryingToLogIn User
	errDecode := json.NewDecoder(r.Body).Decode(&userTryingToLogIn)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", errDecode)
		return
	}

	matchedUserRow := DB.QueryRow("SELECT id, created_at, email FROM users WHERE email = ? AND password = ?;", userTryingToLogIn.Email, userTryingToLogIn.Password)
	fmt.Printf("%+v\n", matchedUserRow)

	w.WriteHeader(http.StatusOK)
}
