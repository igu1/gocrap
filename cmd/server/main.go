package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igu1/gocrap/internal/server"
	"github.com/playwright-community/playwright-go"
)

func main() {
	if err := playwright.Install(); err != nil {
		log.Fatal(err)
		return
	}
	server.RegisterRoutes()
	fmt.Println("Server running at http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
