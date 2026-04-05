package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	cardService services.CardService
}

func NewCardController(cardService services.CardService) *CardController {
	return &CardController{cardService}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to get data", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	if err := c.cardService.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed to create card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type UpdateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req UpdateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	if err := c.cardService.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed update card", err.Error())
	}

	return utils.Success(ctx, "Card update success", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}

	card, err := c.cardService.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card not found", err.Error())
	}

	if err := c.cardService.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Failed to remove card", err.Error())
	}

	return utils.Success(ctx, "Success remove card", publicID)
}

func (c *CardController) GetListCard(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "List ID not valid", err.Error())
	}

	cards, err := c.cardService.GetByListID(listID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to get data cards", err.Error())
	}

	return utils.Success(ctx, "Success get cards", cards)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	card, err := c.cardService.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to get Card", err.Error())
	}

	return utils.Success(ctx, "Success get detail card", card)
}
