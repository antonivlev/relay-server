package frontend

import (
	"html/template"
	"net/http"
	"os"
)

func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/pages/login.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func GetHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/pages/index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func GetFavicon(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("..."))
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	fileData, err := os.ReadFile("frontend/public/scripts/utils.js")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(fileData)
}
