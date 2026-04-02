package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"

	"github.com/gofiber/fiber/v2"
)

type ListController struct {
	listService services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{listService: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Failed to parse data", err.Error())
	}

	if err := c.listService.Create(list); err != nil {
		return utils.BadRequest(ctx, "Failed to created list", err.Error())
	}

	return utils.Success(ctx, "Success created list", list)
}
