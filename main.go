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

	boardRepository := repositories.NewBoardRepository()
	boardMemberRepository := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepository, userRepository, boardMemberRepository)
	boardController := controllers.NewBoardController(boardService)

	listPositionRepository := repositories.NewListPositionRepository()
	listRepository := repositories.NewListRepository()
	listService := services.NewListService(listRepository, boardRepository, listPositionRepository)
	listController := controllers.NewListController(listService)

	cardRepository := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepository, listRepository, userRepository)
	cardController := controllers.NewCardController(cardService)

	routes.Setup(app, userController, boardController, listController, cardController)
	log.Println("Server is running on port: ", port)
	log.Fatal(app.Listen(":" + port))
}
