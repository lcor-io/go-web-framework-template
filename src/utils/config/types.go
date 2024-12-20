package config

import "time"

type RouteType string

const (
	Dynamic RouteType = "dynamic"
	Static  RouteType = "static"

	ConfigFileName      = "config"
	ConfigFileExtension = "toml"
	ConfigDefaultPath   = "."
)

type (
	SubRoutes map[string]Route
	RouteOpts struct {
		Type           RouteType
		Revalidate     time.Duration // Used to set the Cache-Control max-age and stale-while-revalidate directives
		RevalidateTags []string      // Used to invalidate this specific route
	}
	Route struct {
		RouteOpts `mapstructure:",squash"`
		SubRoutes
	}
)

type Config struct {
	Routes map[string]Route
}
