package main

import (
	"campaignservice/db"
	"campaignservice/handler"
	"campaignservice/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	d := db.Connect()
	defer d.Close()
	
	r := chi.NewRouter()

	r.Route("/v1", func(v1 chi.Router) {
		v1.With(utils.MethodGuard("GET")).Get("/delivery", handler.DeliveryHandler)
	})

	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
