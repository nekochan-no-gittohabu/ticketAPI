package main

import (
	"fmt"
	"net/http"
	"os"
	"ticketAPI/api"
	"ticketAPI/ticket"
	"ticketAPI/ticket/store/postgres"

	"github.com/joho/godotenv"
)

const port = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: %w", err)
		return
	}
	store, err := postgres.New(os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.Close()

	service := ticket.New(store)
	server := api.New(service)

	err = http.ListenAndServe(port, server)
	if err != nil {
		fmt.Println(err)
		return
	}
}
