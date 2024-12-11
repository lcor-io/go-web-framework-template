package utils

import (
	"fmt"
	"os"
	"path"
	"slices"
	"time"

	"github.com/gofiber/fiber/v3/log"
)

type routeCache struct {
	Path       string
	Name       string
	ValidUntil time.Time
	Tags       []string
}

type cacheManager struct {
	RouteCachePath string
	Routes         []routeCache
}

func (c *cacheManager) Init() {
	// Create route cache directory
	if err := os.MkdirAll(c.RouteCachePath, 0755); err != nil {
		fmtErr := fmt.Errorf("Error while creating route cache directory %v", err)
		panic(fmtErr)
	}
}

func (c *cacheManager) GetRouteFile(routePath string, validity time.Duration, tags ...string) (*os.File, error) {
	cacheIndex := slices.IndexFunc(c.Routes, func(r routeCache) bool {
		return r.Path == routePath
	})

	// Create the file with the validity if we have a cache miss
	if cacheIndex == -1 {
		log.Info("Cache Miss for route ", routePath)
		f, err := os.CreateTemp(c.RouteCachePath, "route-*.html")
		if err != nil {
			return nil, err
		}
		c.Routes = append(c.Routes, routeCache{
			Path:       routePath,
			Name:       f.Name(),
			ValidUntil: time.Now().Add(validity),
			Tags:       tags,
		})
		return f, nil
	}

	route := c.Routes[cacheIndex]
	f, err := os.Open(route.Name)
	if err != nil {
		return nil, err
	}

	log.Info("Cache Hit for route ", routePath)

	// Cache hit but expired, cleanup the file and restore validity
	if time.Now().After(route.ValidUntil) {
		log.Info("Cache expired")
		if err := os.Truncate(path.Join(c.RouteCachePath, route.Name), 0); err != nil {
			return nil, err
		}
		c.Routes[cacheIndex].ValidUntil = time.Now().Add(validity)
	}

	// Cache hit and cache still valid
	return f, nil
}

func (c *cacheManager) CleanCache() {
	os.RemoveAll(c.RouteCachePath)
}

var CacheManager = cacheManager{
	RouteCachePath: path.Join("cache", "route"),
	Routes:         []routeCache{},
}
