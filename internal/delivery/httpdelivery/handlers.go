package httpdelivery

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type HandlersRegistrar interface {
	Register(router *echo.Echo)
}

type HandlerContainer struct {
	logger *zerolog.Logger
}

func NewHandlers(logger *zerolog.Logger) HandlersRegistrar {
	return &HandlerContainer{logger: logger}
}

func (h *HandlerContainer) Register(router *echo.Echo) {
	router.GET("/:id", h.GetUserByID)
}
