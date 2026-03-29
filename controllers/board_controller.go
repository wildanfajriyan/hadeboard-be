package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	boardService services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{boardService: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Failed read request", err.Error())
	}

	board.OwnerPublicID = userID

	if err := ctx.BodyParser(&board); err != nil {
		return utils.BadRequest(ctx, "Failed read request", err.Error())
	}

	if err := c.boardService.Create(board); err != nil {
		return utils.BadRequest(ctx, "Failed create board", err.Error())
	}

	return utils.Success(ctx, "Success Create Board", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Failed parse data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}

	existingBoard, err := c.boardService.FindByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board not found", err.Error())
	}
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerInternalID = existingBoard.OwnerInternalID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.boardService.Update(board); err != nil {
		return utils.BadRequest(ctx, "Failed update board", err.Error())
	}

	return utils.Success(ctx, "Board updated succesfully", board)
}
