package internal

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/order")
	//g.GET("/:id", handler.GetByID)
	//g.POST("/", handler.Create)
	//g.PUT("/:id", handler.Update)
	g.DELETE("/:id", handler.Delete)
	//g.GET("/list", handler.GetList)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Order not found",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
