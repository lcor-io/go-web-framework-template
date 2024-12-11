package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"

	root "default.app/src/app"
	"default.app/src/utils"
	renderers "default.app/src/utils/renderers"
)

func main() {
	/***
	 * Load environment variables from .env file
	 ***/
	if err := godotenv.Load(); err != nil {
		fmtErr := fmt.Errorf("Error loading .env file: %v", err)
		panic(fmtErr)
	}

	app := fiber.New(fiber.Config{
		AppName: "default.app",
	})

	/***
	 * Setup compression for incoming requests
	 ***/
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	/***
	 * Setup logger in dev
	 ***/
	if os.Getenv("ENV") == "development" {
		app.Use(logger.New(logger.Config{}))
	}

	/***
	 * Setup the cache manager
	 ***/
	utils.CacheManager.Init()

	/***
	 * Serve static files
	 ***/
	if os.Getenv("ENV") == "development" {
		app.Static("/static", "./static", fiber.Static{
			MaxAge: 0,
		})
	} else {
		app.Static("/static", "./static", fiber.Static{
			MaxAge: 3600,
		})
	}

	app.Get("/", func(c fiber.Ctx) error {
		return renderers.StaticRender(&c, root.Index())
	})

	log.Fatal(app.Listen(":42068", fiber.ListenConfig{
		EnablePrintRoutes: os.Getenv("ENV") == "development",
		EnablePrefork:     true,
		OnShutdownSuccess: func() { utils.CacheManager.CleanCache() },
		OnShutdownError:   func(_ error) { utils.CacheManager.CleanCache() },
	}))
}
