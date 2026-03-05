package seed

import (
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"hadeboard-be/utils"
	"log"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin*123")

	admin := models.User{
		Name:     "Super Admin",
		Email:    "admin@example.com",
		Password: password,
		Role:     "admin",
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("Failed to seed admin", err)
	} else {
		log.Println("Admin user seeded")
	}
}
