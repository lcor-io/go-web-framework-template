package renderers

import (
	"context"
	// "os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"default.app/src/utils"
)

type RenderOpts struct {
	revalidate    time.Duration // Used to set the Cache-Control max-age and stale-while-revalidate directives
	revalidateTag string        // Used to invalidate this specific route

	ctx context.Context // Context used to render the component
}

type renderOptFunc func(*RenderOpts)

func defaultOpts() RenderOpts {
	return RenderOpts{
		revalidate:    time.Second * 60 * 60 * 24 * 30, // 30 days
		revalidateTag: "",                              // No tag
		ctx:           context.Background(),
	}
}

func WithContext(ctx context.Context) renderOptFunc {
	return func(o *RenderOpts) {
		o.ctx = ctx
	}
}

func WithRevalidate(revalidate time.Duration) renderOptFunc {
	return func(o *RenderOpts) {
		o.revalidate = revalidate
	}
}

func WithRevalidateTag(tag string) renderOptFunc {
	return func(o *RenderOpts) {
		o.revalidateTag = tag
	}
}

func StaticRender(c fiber.Ctx, component templ.Component, opts ...renderOptFunc) error {
	/***
	* Apply options to the renderer
	***/
	opt := defaultOpts()
	for _, fn := range opts {
		fn(&opt)
	}

	/***
	 * Serve file dynamically in development
	 ***/
	// if os.Getenv("ENV") == "development" {
	// 	return DynamicRender(c, component, opts...)
	// }

	/***
	 * Serve static files in production
	 ***/
	f, err := utils.CacheManager.GetRouteFile(c.Path(), opt.revalidate, opt.revalidateTag)
	if err != nil {
		log.Warnf("Could not create cache file, render component dynamically: %v", err)
		return DynamicRender(&c, component, opts...)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return DynamicRender(&c, component, opts...)
	}

	if stat.Size() == 0 {
		err = component.Render(opt.ctx, f)
		if err != nil {
			log.Warnf("Could not create cache file, render component dynamically: %v", err)
			return DynamicRender(&c, component, opts...)
		}
	}

	cacheAge := time.Now().Sub(stat.ModTime()).Seconds()
	maxAge := opt.revalidate.Seconds()

	c.Set(fiber.HeaderAge, strconv.Itoa(int(cacheAge)))
	c.Set(fiber.HeaderCacheControl, "public, max-age="+strconv.Itoa(int(maxAge)))
	return c.SendFile(f.Name())
}
