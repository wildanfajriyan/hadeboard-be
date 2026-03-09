package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{userService: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Failed parsing data", err.Error())
	}

	if err := c.userService.Register(user); err != nil {
		return utils.BadRequest(ctx, "Failed Register", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Register Success", userResp)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	user, err := c.userService.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Success Login", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResp,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := c.userService.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User Not Found", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "User Found", userResp)
}
