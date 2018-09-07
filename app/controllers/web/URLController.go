package controllers_web

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
}

type URLController struct {
}

func (this *URLController) IndexView(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &User{user})
}

func (this *URLController) LoginView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/login.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
