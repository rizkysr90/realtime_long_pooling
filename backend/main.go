package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"realtime-dashboard-food-delivery/server"
	"realtime-dashboard-food-delivery/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	databaseUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	notificationService := service.NewNotification()
	server := server.NewServer(db, notificationService)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
