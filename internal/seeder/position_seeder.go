package seeder

import (
	"log"

	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// defaultOrgID matches the default organization seeded by
// db/migrations/002_add_organization_id.sql.
const defaultOrgID = "00000000-0000-0000-0000-000000000001"

func SeedPositions(db *gorm.DB) error {
	orgID, err := uuid.FromString(defaultOrgID)
	if err != nil {
		return err
	}

	positions := []model.Position{
		{OrganizationID: &orgID, Name: "Administrator"},
		{OrganizationID: &orgID, Name: "Manager"},
		{OrganizationID: &orgID, Name: "Staff"},
	}

	for _, p := range positions {
		var existing model.Position

		// idempotent: skip if name already exists for this org
		result := db.Where("name = ? AND organization_id = ?", p.Name, orgID).First(&existing)
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
