package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"default.app/cmd/app/builder"
)

func main() {

	routeTree := buildRouteTree(builder.BASE_PATH)
	if routeTree == nil {
		panic("route tree should not be nil")
	}

	// Build the entrypoint
	f, e := os.Create("src/app/entrypoint_app.go")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	tmpl, tmplErr := template.ParseFiles("cmd/app/templates/entrypoint.tmpl")
	if tmplErr != nil {
		panic(tmplErr)
	}
	entrypointErr := tmpl.Execute(f, routeTree)
	if entrypointErr != nil {
		panic(entrypointErr)
	}

}

func buildRouteTree(path string) *builder.Route {

	opts := new(builder.RouteOpts)
	routes := []string{}
	subRoutes := []*builder.Route{}

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, e := range entries {

		filePath := filepath.Join(path, e.Name())
		fileExt := filepath.Ext(filePath)
		fileName := strings.ReplaceAll(e.Name(), fileExt, "")

		if !e.IsDir() && fileExt != ".templ" {
			continue
		}

		if fileName == "index" && fileExt == ".templ" {
			opts = parseOptsFromFile(filePath)
			continue
		}

		if fileName != "index" && fileExt == ".templ" {
			routes = append(routes, fileName)
			continue
		}

		if subRoute := buildRouteTree(filePath); subRoute != nil {
			subRoutes = append(subRoutes, subRoute)
		}
	}

	if len(routes) == 0 && len(subRoutes) == 0 && opts == nil {
		return nil
	}

	route := builder.Route{
		Path:      strings.ReplaceAll(path, builder.BASE_PATH, ""),
		Opts:      opts,
		Routes:    routes,
		SubRoutes: subRoutes,
	}

	if path != builder.BASE_PATH {
		if err := buildRoute(path, route); err != nil {
			panic(err)
		}
	}

	return &route
}

func buildRoute(filePath string, route builder.Route) error {

	f, e := os.Create(filePath + "/index_app.go")
	if e != nil {
		return e
	}
	defer f.Close()

	tmpl, tmplErr := template.ParseFiles("cmd/app/templates/index.tmpl")
	if tmplErr != nil {
		return tmplErr
	}

	if indexErr := tmpl.Execute(f, route); indexErr != nil {
		return indexErr
	}

	return nil
}

func parseOptsFromFile(filePath string) *builder.RouteOpts {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	opts := builder.RouteOpts{
		Package:        "",
		Type:           builder.Static,
		Revalidate:     time.Second * 60 * 60 * 24 * 30, // 30 days,
		RevalidateTags: []string{},                      // No tags
	}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, "package") {
			filePackage := strings.Replace(line, "package ", "", 1)
			opts.Package = filePackage
			continue
		}
		if strings.Contains(line, "var revalidateTags =") {
			opts.RevalidateTags = getRouteRevalidateTags(line)
			continue
		}
		if strings.Contains(line, "var revalidate =") {
			opts.Revalidate = getRouteRevalidate(line)
			if opts.Revalidate == 0 {
				opts.Type = builder.Dynamic
			}
		}
	}

	return &opts
}

func getRouteRevalidate(line string) time.Duration {
	rawDuration := strings.Split(line, "=")[1]
	rawDuration = strings.TrimSpace(rawDuration)
	duration, err := strconv.ParseInt(rawDuration, 10, 64)
	if err != nil {
		panic(err)
	}

	if duration <= 0 {
		return 0
	}
	return time.Second * time.Duration(duration)
}

func getRouteRevalidateTags(line string) []string {
	rawTagLine := strings.Split(line, "=")[1]
	rawTagLine = strings.Replace(rawTagLine, "[]string{", "", 1)
	rawTagLine = strings.Replace(rawTagLine, "}", "", 1)
	rawTags := strings.Split(rawTagLine, ",")

	tags := []string{}

	for _, tag := range rawTags {
		tag = strings.TrimSpace(tag)
		tag = strings.ReplaceAll(tag, "\"", "")
		tags = append(tags, tag)
	}

	return tags
}
