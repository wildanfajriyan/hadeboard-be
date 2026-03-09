package routes

import (
	"hadeboard-be/config"
	"hadeboard-be/controllers"
	"hadeboard-be/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
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

	apiProtected := app.Group("/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	userGroup := apiProtected.Group("/users")
	userGroup.Get("/:id", userController.GetUser)
}
