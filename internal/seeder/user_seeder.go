package seeder

import (
	"log"

	"Microservice/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []struct {
		EmployeeID   string
		Email        string
		Password     string
		Role         int
		FirstName    string
		LastName     string
		Phone        string
		Access       bool
		PositionName string // used to resolve PositionID FK
	}{
		{
			EmployeeID:   "EMP001",
			Email:        "admin@approval.com",
			Password:     "Test1234!",
			Role:         99, // 99 = admin
			FirstName:    "Super",
			LastName:     "Admin",
			Phone:        "08123456789",
			Access:       true,
			PositionName: "Administrator",
		},
		{
			EmployeeID:   "EMP002",
			Email:        "manager1@approval.com",
			Password:     "Test1234!",
			Role:         1, // 1 = user
			FirstName:    "John",
			LastName:     "Doe",
			Phone:        "08234567890",
			Access:       true,
			PositionName: "Manager",
		},
		{
			EmployeeID:   "EMP002",
			Email:        "manager2@approval.com",
			Password:     "Test1234!",
			Role:         1, // 1 = user
			FirstName:    "Steve",
			LastName:     "Wozniak",
			Phone:        "08234567890",
			Access:       true,
			PositionName: "Manager",
		},
		{
			EmployeeID:   "EMP003",
			Email:        "staff1@approval.com",
			Password:     "Test1234!",
			Role:         1,
			FirstName:    "Jane",
			LastName:     "Smith",
			Phone:        "08345678901",
			Access:       true,
			PositionName: "Staff",
		},
		{
			EmployeeID:   "EMP004",
			Email:        "staff2@approval.com",
			Password:     "Test1234!",
			Role:         1,
			FirstName:    "Grace",
			LastName:     "Mary",
			Phone:        "08345678901",
			Access:       true,
			PositionName: "Staff",
		},
	}

	for _, u := range users {
		var existing model.User

		// 1. skip if email already exists
		result := db.Where("email = ?", u.Email).First(&existing)
		if result.Error == nil {
			log.Printf("⏭  User '%s' already exists, skipping\n", u.Email)
			continue
		}

		// 2. resolve PositionID FK — query from positions table
		var position model.Position
		if err := db.Where("name = ?", u.PositionName).First(&position).Error; err != nil {
			log.Printf("❌ Position '%s' not found for user '%s'\n", u.PositionName, u.Email)
			return err
		}

		// 3. hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 4. build and insert user
		newUser := model.User{
			PositionID: position.ID, // resolved FK from step 2
			EmployeeID: u.EmployeeID,
			Email:      u.Email,
			Password:   string(hash),
			Role:       u.Role,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Phone:      u.Phone,
			Access:     u.Access,
		}

		if err := db.Create(&newUser).Error; err != nil {
			return err
		}
		log.Printf("✅ Seeded user: %s %s (%s)\n", u.FirstName, u.LastName, u.Email)
	}

	return nil
}
