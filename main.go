package main

import (
	"campaignservice/db"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	d := db.Connect()
	defer d.Close()

}
