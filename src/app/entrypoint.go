package app

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/rewrite"

	"default.app/src/app/page-1"
	"default.app/src/app/page-2"
)

func RegisterRoutes(app *fiber.App) {
	// Middleware to check if the request is coming from HTMX, and is eventually boosted.
	// Store the result in the context locals
	app.Use(func(c fiber.Ctx) error {
		isHTMXRequest, errHtmx := strconv.ParseBool(c.Get("hx-request"))
		isBoosted, errBoosted := strconv.ParseBool(c.Get("hx-boosted"))

		c.Locals("isHTMXRequest", errHtmx == nil && isHTMXRequest)
		c.Locals("isBoosted", errBoosted == nil && isBoosted)

		return c.Next()
	})

	// Default to page 1
	app.Use(rewrite.New(rewrite.Config{
		Rules: map[string]string{
			"/": "/page-1",
		},
	}))

	files, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		println(file.Name())
	}

	// Register every sub-routes
	page1.RegisterRoute(app.Group("/page-1"))
	page2.RegisterRoute(app.Group("/page-2"))
}
