package main

import (
	"fmt"
	"net/http"

	templates "github.com/Sheriff-Hoti/go-url-shortener/templates"
	templ "github.com/a-h/templ"
)

func main() {
	component := templates.Page()
	http.Handle("/", templ.Handler(component))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
