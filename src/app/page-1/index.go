package page1

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"

	"default.app/src/components"
	"default.app/src/utils/renderers"
)

func RegisterRoute(router fiber.Router) {
	router.Get("/", func(ctx fiber.Ctx) error {
		isBoosted, ok := ctx.Locals("isBoosted").(bool)

		// By leveraging the boot feature in HTMX, we can return only the content of the page and reduce the payload
		if ok && isBoosted {
			return renderers.DynamicRender(&ctx, Index())
		}

		// But if JS is not enabled, we can still render the full page
		childContext := templ.WithChildren(ctx.Context(), Index())
		return renderers.DynamicRender(&ctx, components.MainLayout(), renderers.WithContext(childContext))
	})

	router.Get("/tab1", func(ctx fiber.Ctx) error {
		return renderers.DynamicRender(&ctx, Tab(1, "Content of tab 1"))
	})
	router.Get("/tab2", func(ctx fiber.Ctx) error {
		return renderers.DynamicRender(&ctx, Tab(2, "Content of tab 2"))
	})
	router.Get("/tab3", func(ctx fiber.Ctx) error {
		return renderers.DynamicRender(&ctx, Tab(3, "Content of tab 3"))
	})
}
