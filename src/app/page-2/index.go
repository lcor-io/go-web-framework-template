package page2

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"

	"default.app/src/components"
	"default.app/src/utils/renderers"
)

func RegisterRoute(router fiber.Router) {
	router.Get("/", func(ctx fiber.Ctx) error {
		childContext := templ.WithChildren(ctx.Context(), Index())
		return renderers.StaticRender(ctx, components.MainLayout(), renderers.WithContext(childContext))
	})
}
