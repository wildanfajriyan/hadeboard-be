package routes

import (
	"hadeboard-be/config"
	"hadeboard-be/controllers"
	"hadeboard-be/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App,
	userController *controllers.UserController,
	boardController *controllers.BoardController,
	listController *controllers.ListController,
	cardController *controllers.CardController,
) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	app.Use(logger.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	app.Post("/v1/auth/register", userController.Register)
	app.Post("/v1/auth/login", userController.Login)

	apiProtected := app.Group("/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	userGroup := apiProtected.Group("/users")
	userGroup.Get("/page", userController.GetUserPagination)
	userGroup.Get("/:id", userController.GetUser)
	userGroup.Put("/:id", userController.UpdateUser)
	userGroup.Delete("/:id", userController.DeleteUser)

	boardGroup := apiProtected.Group("/boards")
	boardGroup.Get("/my", boardController.GetMyBoardPaginate)
	boardGroup.Post("/", boardController.CreateBoard)
	boardGroup.Put("/:id", boardController.UpdateBoard)
	boardGroup.Post("/:id/members", boardController.AddBoardMembers)
	boardGroup.Delete("/:id/members", boardController.RemoveBoardMembers)
	boardGroup.Get("/:board_id/lists", listController.GetListOnBoard)

	listGroup := apiProtected.Group("/lists")
	listGroup.Post("/", listController.CreateList)
	listGroup.Put("/:id", listController.UpdateList)
	listGroup.Delete("/:id", listController.DeleteList)
	listGroup.Get("/:list_id/cards", cardController.GetListCard)

	cardGroup := apiProtected.Group("/cards")
	cardGroup.Post("/", cardController.CreateCard)
	cardGroup.Put("/:id", cardController.UpdateCard)
	cardGroup.Delete("/:id", cardController.DeleteCard)
	cardGroup.Get("/:id", cardController.GetCardDetail)
}
