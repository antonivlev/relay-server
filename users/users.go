package users

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	NumberOfTokens float64   `json:"numberOfTokens"`
}

var DB *sql.DB

func GetUser(w http.ResponseWriter, r *http.Request) {
	email, _, _ := r.BasicAuth()
	rowUser := DB.QueryRow("SELECT id, created_at, email, number_of_tokens FROM users WHERE email = ?;", email)

	var user User
	errScan := rowUser.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.NumberOfTokens)
	if errScan != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errScan)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rowsUsers, errQuery := DB.Query("SELECT id, created_at, email, number_of_tokens FROM users;")
	if errQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", errQuery)
		return
	}

	var users []User
	for rowsUsers.Next() {
		var user User
		errScan := rowsUsers.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.NumberOfTokens)
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
	var emailAndPassword struct{ Email, Password string }
	errDecode := json.NewDecoder(r.Body).Decode(&emailAndPassword)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", errDecode)
		return
	}

	matchedUserRow := DB.QueryRow("SELECT id, created_at, email FROM users WHERE email = ? AND password = ?;", emailAndPassword.Email, emailAndPassword.Password)
	var matchedUser User
	errScan := matchedUserRow.Scan(&matchedUser.ID, &matchedUser.CreatedAt, &matchedUser.Email)
	if errScan != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%v", "Email and password do not match.")
		return
	}

	emailColonPassword := emailAndPassword.Email + ":" + emailAndPassword.Password
	accessToken := base64.StdEncoding.EncodeToString([]byte(emailColonPassword))

	var accessTokenResponse struct {
		AccessToken string `json:"accessToken"`
	}
	accessTokenResponse.AccessToken = accessToken

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessTokenResponse)
}
