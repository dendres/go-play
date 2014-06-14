package main

import (
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	h := "<!DOCTYPE html><body>\n"
	problem_name := ""
	path := ""

	files, err := ioutil.ReadDir(".")
	if err == nil {
		for _, fi := range files {
			if fi.IsDir() == true {
				problem_name = fi.Name()
				path = "/" + problem_name + "/" + problem_name + ".svg"
				h += "<object data=\"" + path + "\" width=\"100%\" height=\"800\" type=\"image/svg+xml\">" + problem_name + "</object>"
			}
		}
	}

	m := martini.Classic()
	m.Use(martini.Static("."))
	m.Get("/", func() string {
		return h
	})

	log.Fatal(http.ListenAndServe(":8080", m))
}
