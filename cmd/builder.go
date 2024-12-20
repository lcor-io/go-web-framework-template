package main

import (
	"bufio"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
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

type Conf struct {
	Routes map[string]Route
}

func main() {
	const fileConfigurationName = "conf"
	const fileConfigurationType = "toml"
	const defaultConfPath = "."

	viper.SetConfigName(fileConfigurationName)
	viper.SetConfigType(fileConfigurationType)
	viper.AddConfigPath(defaultConfPath)

	const basePath = "src/app/"

	filepath.WalkDir(basePath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			panic(e)
		}

		// Ignore files that are not .templ
		if filepath.Ext(d.Name()) != ".templ" {
			return nil
		}

		filePath := strings.Replace(s, basePath, "", 1)
		filePath = strings.ReplaceAll(filePath, ".templ", "")
		filePath = strings.ReplaceAll(filePath, "/index", "")
		filePathExploded := strings.Split(filePath, "/")

		fileOpts := parseOptsFromFile(s)

		routePath := "routes." + strings.Join(filePathExploded, ".")
		viper.Set(routePath, filePath)
		viper.Set(routePath+".type", fileOpts.Type)
		viper.Set(routePath+".revalidate", fileOpts.Revalidate)
		viper.Set(routePath+".revalidateTags", fileOpts.RevalidateTags)

		return nil
	})

	if err := viper.WriteConfigAs(path.Join(defaultConfPath, fileConfigurationName+"."+fileConfigurationType)); err != nil {
		panic(err)
	}
}

func parseOptsFromFile(filePath string) RouteOpts {
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
		RevalidateTags: []string{},                      // No tags
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
		tag = strings.ReplaceAll(tag, "\"", "")
		tags = append(tags, tag)
	}

	return tags
}
