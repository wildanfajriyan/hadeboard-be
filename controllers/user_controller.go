package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return utils.BadRequest(ctx, err.Error(), err.Error())
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

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.userService.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to Get Data", err.Error())
	}

	var userResp []models.UserResponse
	err = copier.Copy(&userResp, &users)
	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		utils.NotFoundPagination(ctx, "Users not found", userResp, meta)
	}

	return utils.SuccessPagination(ctx, "Users Found", userResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Failed parse data", err.Error())
	}
	user.PublicID = publicID

	if err := c.userService.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Failed parse data", err.Error())
	}

	userUpdated, err := c.userService.GetByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed get data", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Error parsing data", err.Error())
	}

	return utils.Success(ctx, "Success Update Data", userResp)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.userService.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Failed remove data", err.Error())
	}

	return utils.Success(ctx, "Success Remove Data", id)
}
