package utils

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"default.app/src/utils"
)

type StaticRenderOpts struct {
	revalidate    time.Duration // Used to set the Cache-Control max-age and stale-while-revalidate directives
	revalidateTag string        // Used to invalidate this specific route

	templHandlers []func(*templ.ComponentHandler) // Handlers to apply to the component
}

type staticRenderOptFunc func(*StaticRenderOpts)

func defaultOpts() StaticRenderOpts {
	return StaticRenderOpts{
		revalidate:    time.Second * 60 * 60 * 24 * 30,   // 30 days
		revalidateTag: "",                                // No tag
		templHandlers: []func(*templ.ComponentHandler){}, // No handlers
	}
}

func WithRevalidate(revalidate time.Duration) staticRenderOptFunc {
	return func(o *StaticRenderOpts) {
		o.revalidate = revalidate
	}
}

func WithRevalidateTag(tag string) staticRenderOptFunc {
	return func(o *StaticRenderOpts) {
		o.revalidateTag = tag
	}
}

func withTemplHandlers(handlers ...func(*templ.ComponentHandler)) staticRenderOptFunc {
	return func(o *StaticRenderOpts) {
		o.templHandlers = append(o.templHandlers, handlers...)
	}
}

func StaticRender(c *fiber.Ctx, component templ.Component, opts ...staticRenderOptFunc) error {
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
	if os.Getenv("ENV") == "development" {
		return DynamicRender(c, component, opt.templHandlers...)
	}

	/***
	 * Serve static files in production
	 ***/
	f, err := utils.CacheManager.GetRouteFile((*c).Path(), opt.revalidate, opt.revalidateTag)
	if err != nil {
		log.Warnf("Could not create cache file, render component dynamically: %v", err)
		return DynamicRender(c, component, opt.templHandlers...)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return DynamicRender(c, component, opt.templHandlers...)
	}

	if stat.Size() == 0 {
		err = component.Render(context.Background(), f)
		if err != nil {
			log.Warnf("Could not create cache file, render component dynamically: %v", err)
			return DynamicRender(c, component, opt.templHandlers...)
		}
	}

	cacheAge := time.Now().Sub(stat.ModTime()).Seconds()
	maxAge := opt.revalidate.Seconds()

	(*c).Set(fiber.HeaderAge, strconv.Itoa(int(cacheAge)))
	(*c).Set(fiber.HeaderCacheControl, "public, max-age="+strconv.Itoa(int(maxAge)))
	return (*c).SendFile(f.Name())
}
