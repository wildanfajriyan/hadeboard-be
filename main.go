package main

import (
	"hadeboard-be/config"
	"hadeboard-be/controllers"
	"hadeboard-be/database/seed"
	"hadeboard-be/repositories"
	"hadeboard-be/routes"
	"hadeboard-be/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	seed.SeedAdmin()
	port := config.AppConfig.AppPort

	app := fiber.New()

	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, userController)
	log.Println("Server is running on port: ", port)
	log.Fatal(app.Listen(":" + port))
}
