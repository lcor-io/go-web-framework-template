package page2

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"

	"default.app/src/components"
	renderers "default.app/src/utils/renderers"
)

func RegisterRoute(router fiber.Router) {
	router.Get("/", func(ctx fiber.Ctx) error {
		isHXRequest := ctx.Get("hx-request")
		isBoosted := ctx.Get("hx-boosted")

		println("is htmx request:", isHXRequest)
		println("is boosted:", isBoosted)

		childContext := templ.WithChildren(ctx.Context(), Index())
		return renderers.DynamicRender(&ctx, components.MainLayout("Page 2"), renderers.WithContext(childContext))
	})
}
