package main

import (
	"net/http"
	"os"
)

type error struct {
	Error string
}

func main() {
	http.HandleFunc("/match", match)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "http"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
