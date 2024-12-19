package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type RouteType string

const (
	Dynamic RouteType = "dynamic"
	Static  RouteType = "static"
)

type RouteOpts struct {
	Type           RouteType
	Revalidate     time.Duration // Used to set the Cache-Control max-age and stale-while-revalidate directives
	RevalidateTags []string      // Used to invalidate this specific route
}
type Route struct {
	Path []string
	Opts RouteOpts
}

var files []Route

func main() {
	// Create and truncate the conf file
	_, err := os.Create("conf.toml")
	if err != nil {
		panic(err)
	}
	if err := os.Truncate("conf.toml", 0); err != nil {
		panic(err)
	}

	path := path.Join("src", "app")
	walkDir(path)

	// Write the routes to the conf file
	f, err := os.OpenFile("conf.toml", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("# This file is auto-generated - DO NOT EDIT\n\n")
	for _, file := range files {
		_, err := f.WriteString(fmt.Sprintf("[%s]\ntype = \"%s\"\nrevalidate = %d\nrevalidateTags = %v\n\n", strings.Join(file.Path, "."), file.Opts.Type, file.Opts.Revalidate, file.Opts.RevalidateTags))
		if err != nil {
			panic(err)
		}
	}
}

func walkDir(dirPath string) {
	filepath.WalkDir(dirPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			panic(e)
		}

		// Add files if they have the .templ extension
		if filepath.Ext(d.Name()) == ".templ" {

			filePath := strings.Replace(s, path.Join("src", "app"), "", 1)
			filePath = strings.Replace(s, ".templ", "", 1)
			filePathExploded := strings.Split(filePath, "/")
			filePathExploded = filePathExploded[2:]

			fileOpts := getOptsFromFile(s)

			files = append(files, Route{
				Path: filePathExploded,
				Opts: fileOpts,
			})
		}

		return nil
	})
}

func getOptsFromFile(filePath string) RouteOpts {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	opts := RouteOpts{
		Type:           Static,
		Revalidate:     time.Second * 60 * 60 * 24 * 30, // 30 days,
		RevalidateTags: []string{},
	}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, "var revalidate =") {
			opts.Revalidate = getRouteRevalidate(line)
			if opts.Revalidate == 0 {
				opts.Type = Dynamic
			}
		}
		if strings.Contains(line, "var revalidateTags =") {
			opts.RevalidateTags = getRouteRevalidateTags(line)
		}
	}

	return opts
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
		tags = append(tags, tag)
	}

	return tags
}
