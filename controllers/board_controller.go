package controllers

import (
	"hadeboard-be/internal/models"
	"hadeboard-be/services"
	"hadeboard-be/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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

func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed parsing data", err.Error())
	}

	if err := c.boardService.AddMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed added members", err.Error())
	}

	return utils.Success(ctx, "Member added succesfully", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed parsing data", err.Error())
	}

	if err := c.boardService.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed removed members", err.Error())
	}

	return utils.Success(ctx, "Member removed succesfully", nil)
}

func (c *BoardController) GetMyBoardPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userID, err := uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Failed read request", err.Error())
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.boardService.GetMyBoardPaginate(userID.String(), filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to Get Data", err.Error())
	}

	var boardResp []models.BoardResponse
	err = copier.Copy(&boardResp, &boards)
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
		utils.NotFoundPagination(ctx, "Boards not found", boardResp, meta)
	}

	return utils.SuccessPagination(ctx, "Boards Found", boardResp, meta)
}
