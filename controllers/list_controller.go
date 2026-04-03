package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Failed to parse data", err.Error())
	}

	existingList, err := c.listService.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.listService.Update(list); err != nil {
		return utils.BadRequest(ctx, "Failed updated list", err.Error())
	}

	updatedList, err := c.listService.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found after updated", err.Error())
	}

	return utils.Success(ctx, "Success updated list", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}

	lists, err := c.listService.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	return utils.Success(ctx, "Success get lists", lists)
}
