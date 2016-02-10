package driver

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"serve4/serve"
	"strings"

	"github.com/kirubasankars/metal"
)

type FileSystem struct{}

func (fs *FileSystem) GetPath(server serve.Server, url string) (string, []string, string, int, error) {
	url = strings.TrimRight(url, "/")
	if url == "" {
		return "", nil, "", 0, errors.New("Empty Url")
	}

	parts := strings.Split(url, "/")

	appName := ""
	appPath := ""
	sitePath := ""
	length := 0
	sites := make([]string, 0, 0)
	moduleName := ""

	for idx, path := range parts {
		if idx == 0 {
			continue
		}

		if idx == 1 {
			appPath = filepath.Join(server.Path(), "apps", path)
			sitePath = appPath
			if f, _ := os.Stat(appPath); f != nil {
				appName = path
				length += 1 + len(appName)
			} else {
				modulePath := filepath.Join(server.Path(), "modules", path)
				if f, _ := os.Stat(modulePath); f != nil {
					moduleName = path
					length += 1 + len(moduleName)
					return "", nil, moduleName, length, nil
				} else {
					return "", nil, "", 0, errors.New(path + " is not a application/module")
				}
			}
			continue
		}

		sitePath = filepath.Join(sitePath, path)
		if f, _ := os.Stat(sitePath); f != nil {
			sites = append(sites, path)
			length += 1 + len(path)
		} else {
			modulePath := filepath.Join(server.Path(), "modules", path)
			if f, _ := os.Stat(modulePath); f != nil {
				moduleName = path
				length += 1 + len(moduleName)
				return "", nil, moduleName, length, nil
			} else {
				return appName, sites, "", length, errors.New(path + " is not a site/module")
			}
			break
		}
	}

	return appName, sites, moduleName, length, nil
}

func (fs *FileSystem) GetApplication(ctx serve.ServeContext, name string) *serve.Site {
	app := new(serve.Site)
	app.Name = name
	app.URI = "/" + name
	app.Path = filepath.Join(ctx.Server.Path(), "apps", name)
	app.Config = getConfig(app.Path)
	return app
}

func (fs *FileSystem) GetSite(ctx serve.ServeContext, parent *serve.Site, name string) *serve.Site {
	site := new(serve.Site)
	site.Name = name
	site.Path = filepath.Join(parent.Path, name)
	site.Config = getConfig(site.Path)
	return site
}

func (fs *FileSystem) GetModule(ctx serve.ServeContext, name string) *serve.Module {
	module := new(serve.Module)
	module.Name = name
	module.URI = "/" + name
	module.Path = filepath.Join(ctx.Server.Path(), "modules", name)
	module.Config = getConfig(module.Path)

	if f, _ := os.Stat(filepath.Join(module.Path, "app")); f != nil {
		module.AppEnabled = true
	}
	if f, _ := os.Stat(filepath.Join(module.Path, "_auth")); f != nil {
		module.AuthEnabled = true
	}

	return module
}

func getConfig(path string) *metal.Metal {
	configFile := filepath.Join(path, "config.json")
	byt, err := ioutil.ReadFile(configFile)
	if err == nil {
		m := metal.NewMetal()
		m.Parse(byt)
		return m
	}
	return nil
}
