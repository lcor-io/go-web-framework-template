package builder

import "time"

type RouteType string

const (
	Dynamic RouteType = "dynamic"
	Static  RouteType = "static"
)

const BASE_PATH = "src/app/"

type RouteOpts struct {
	Package        string
	Type           RouteType
	Revalidate     time.Duration // Used to set the Cache-Control max-age and stale-while-revalidate directives
	RevalidateTags []string      // Used to invalidate this specific route
}

type Route struct {
	Path string
	Opts *RouteOpts

	Routes    []string
	SubRoutes []*Route
}
