package main

import (
	"fmt"
	"github.com/kidstuff/toys/view"
	"net/http"
)

func main() {
	tmpl := view.NewView("templates")
	tmpl.ResourcePrefix = "/statics"
	tmpl.Watch = true
	if err := tmpl.Parse("default"); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Loading default template")

	http.Handle("/statics/",
		http.StripPrefix("/statics/",
			http.FileServer(http.Dir("statics"))))
	fmt.Println("Handle static file at /statics")

	http.Handle("/", Handler("/", tmpl))
	fmt.Println("start server")

	http.ListenAndServe(":8080", nil)
}
