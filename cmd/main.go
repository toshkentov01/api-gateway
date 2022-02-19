package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/toshkentov01/alif-tech-task/api-gateway/api/docs" //register swagger
	"github.com/toshkentov01/alif-tech-task/api-gateway/config"
	"github.com/toshkentov01/alif-tech-task/api-gateway/pkg/middleware"
	"github.com/toshkentov01/alif-tech-task/api-gateway/pkg/routes"
	"github.com/toshkentov01/alif-tech-task/api-gateway/pkg/utils"
)

var (
	fiberConfig = config.FiberConfig()
	appConfig = config.Config()
)

// @title Alif Tech Task's API
// @description This is an auto-generated API Docs for Alif Tech's Task.
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	app := fiber.New(fiberConfig)

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	middleware.FiberMiddleware(app)

	jwtRoleAuthorizer, err := middleware.NewJWTRoleAuthorizer(appConfig)
	if err != nil {
		log.Fatal("Could not initialize JWT Role Authorizer")
	}

	app.Use(middleware.NewAuthorizer(jwtRoleAuthorizer))
	routes.SwaggerRoute(app)
	routes.UserRoutes(app)

	// Start server (with or without graceful shutdown).
	if config.Config().Environment == "develop" {
		utils.StartServer(app)

	} else {
		utils.StartServerWithGracefulShutdown(app)
	}

}
