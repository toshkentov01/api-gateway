package errors

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	// InternalMsg - error message for internal server error
	InternalMsg = "Internal Server Error"
	// NotEnoughRights - error message for lack of rights
	NotEnoughRights = "Not Enough Rights"
)

// ErrorResponse is a customized error type
type ErrorResponse struct {
	Code    int
	Message string
}

//AbortWithBadRequest handles bad request error
func AbortWithBadRequest(c *fiber.Ctx, err error, msg string) bool {
	if err == nil {
		return false
	}
	c.Status(http.StatusBadRequest)
	err = c.JSON(ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
	return true
}

//AbortWithInternal handles internal error
func AbortWithInternal(c *fiber.Ctx, err error, msg string) bool {
	if err == nil {
		return false
	}
	c.Status(http.StatusInternalServerError)
	err = c.JSON(ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
	return true
}

//AbortWithUnauthorized handles unauthorized user error
func AbortWithUnauthorized(c *fiber.Ctx, err error, msg string) bool {
	if err == nil {
		return false
	}
	c.Status(http.StatusUnauthorized)
	err = c.JSON(ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: msg,
	})
	return true
}
