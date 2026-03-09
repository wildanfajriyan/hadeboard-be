package routes

import (
	"hadeboard-be/controllers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, userController *controllers.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	app.Use(logger.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	app.Post("/v1/auth/register", userController.Register)
	app.Post("/v1/auth/login", userController.Login)
}
