package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/match", match)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
