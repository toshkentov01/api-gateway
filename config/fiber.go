package config

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: 	time.Second * time.Duration(Config().ServerReadTimeout),
		BodyLimit:		1 << 27, // 100 Mb
	}
}
