package routes

import (
	"hadeboard-be/controllers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, userController *controllers.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	app.Post("/v1/auth/register", userController.Register)
}
