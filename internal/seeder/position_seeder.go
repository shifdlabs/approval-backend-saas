package seeder

import (
	"log"

	"Microservice/model"

	"gorm.io/gorm"
)

func SeedPositions(db *gorm.DB) error {
	positions := []model.Position{
		{Name: "Administrator"},
		{Name: "Manager"},
		{Name: "Staff"},
	}

	for _, p := range positions {
		var existing model.Position

		// idempotent: skip if name already exists
		result := db.Where("name = ?", p.Name).First(&existing)
		if result.Error == nil {
			log.Printf("⏭  Position '%s' already exists, skipping\n", p.Name)
			continue
		}

		if err := db.Create(&p).Error; err != nil {
			return err
		}
		log.Printf("✅ Seeded position: %s\n", p.Name)
	}

	return nil
}
