package app

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/rewrite"

	"default.app/src/app/page-1"
	"default.app/src/app/page-2"
)

func RegisterRoutes(app *fiber.App) {
	// Middleware to check if the request is an htmx request, and is eventually boosted
	app.Use(func(c fiber.Ctx) error {
		isHTMXRequest, _ := strconv.ParseBool(c.Get("hx-request"))
		isBoosted, _ := strconv.ParseBool(c.Get("hx-boosted"))

		c.Locals("isHTMXRequest", isHTMXRequest)
		c.Locals("isBoosted", isBoosted)

		return c.Next()
	})

	// Default to page 1
	app.Use(rewrite.New(rewrite.Config{
		Rules: map[string]string{
			"/": "/page-1",
		},
	}))

	// Register every sub-routes
	page1.RegisterRoute(app.Group("/page-1"))
	page2.RegisterRoute(app.Group("/page-2"))
}
