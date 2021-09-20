package httpdelivery

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getUserByIDRequest struct {
	UserID int `param:"id"`
}

func (h *HandlerContainer) GetUserByID(c echo.Context) error {
	requestParams := getUserByIDRequest{}
	var userID int64
	echo.PathParamsBinder(c).Int64("id", &userID)

	if err := c.Bind(&requestParams); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid params", "verbose": err.Error()})
	}

	fmt.Println(userID)
	return nil
}
