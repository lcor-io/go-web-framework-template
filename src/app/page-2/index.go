package page2

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
			return renderers.StaticRender(ctx, Index())
		}

		// But if JS is not enabled, we can still render the full page
		childContext := templ.WithChildren(ctx.Context(), Index())
		return renderers.StaticRender(ctx, components.MainLayout(), renderers.WithContext(childContext))
	})
}
