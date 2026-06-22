package main

import (
	"log"

	"Microservice/config"
	"Microservice/internal/seeder"
)

func main() {
	// same as your main.go — load config
	envConf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables!\n", err.Error())
	}

	// same as your main.go — connect to DB
	db := config.DatabaseConnection(&envConf)

	log.Println("🌱 Starting seeder...")

	// run all seeders
	if err := seeder.Run(db); err != nil {
		log.Fatal("❌ Seeding failed:", err)
	}

	log.Println("🎉 Database seeded successfully!")
}
