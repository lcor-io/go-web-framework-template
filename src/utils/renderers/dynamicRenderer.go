package renderers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func DynamicRender(c *fiber.Ctx, component templ.Component, opts ...renderOptFunc) error {
	/***
	* Apply options to the renderer
	***/
	opt := defaultOpts()
	for _, fn := range opts {
		fn(&opt)
	}

	(*c).Set("Content-Type", "text/html")
	return component.Render(opt.ctx, ((*c).Response().BodyWriter()))
}
