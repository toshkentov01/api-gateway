package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/toshkentov01/alif-tech-task/api-gateway/api/controllers"
)

// UserRoutes func for describe group of public routes.
func UserRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api")

	// Routes For POST Method:
	route.Post("/create-identified-user/", controllers.SignUpFully)
	route.Post("/create-unidentified-user/", controllers.CreateUnidentifiedUser)
	route.Post("/user/income/", controllers.Income)
	route.Post("/user/expense/", controllers.Expense)

	// Routes For GET Method:
	route.Get("/check-user-account/", controllers.CheckUserAccount)
	route.Get("/user/balance/", controllers.GetBalance)
	route.Get("/user/operations/", controllers.ListOperationsByType)

}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
