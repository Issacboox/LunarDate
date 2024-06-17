package handler

import (
	p "bam/internal/core/port"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type LunarDateHandler struct {
	lunarDateService p.LunarDateService
}

func NewLunarDateHandler(lunarDateService p.LunarDateService) *LunarDateHandler {
	return &LunarDateHandler{lunarDateService}
}

func (h *LunarDateHandler) GetLunarDate(c *fiber.Ctx) error {
	date := c.Params("date")
	lunarDate, err := h.lunarDateService.GetLunarDate(date)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"date":   lunarDate,
	})
}
