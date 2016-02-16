package serve

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"
)

type ServeContext struct {
	Server      *Server
	Application *Site
	Site        *Site
	Module      *Module

	URILength int
}

func NewContext(server *Server, url string) *ServeContext {
	ctx := new(ServeContext)
	ctx.Server = server

	system := server.System
	appName, siteNames, moduleName, length, err := system.GetPath(*server, url)
	if err != nil {
		fmt.Println(err)
	}
	ctx.URILength = length

	getApp(ctx, appName)
	getSite(ctx, siteNames)
	getModule(ctx, moduleName)

	ctx.Server = server

	return ctx
}

func getApp(ctx *ServeContext, name string) {
	server := ctx.Server
	if name == "" {
		return
	}
	if app, p := server.Apps[name]; p == false {
		app := server.System.GetApplication(*ctx, name)
		app.jar = createJar(app.Config)
		app.Sites = make(map[string]*Site)
		server.Apps[name] = app
		ctx.Application = app
	} else {
		ctx.Application = app
	}
}

func getSite(ctx *ServeContext, sites []string) {
	server := ctx.Server
	app := ctx.Application

	if !(len(sites) == 0 || app == nil) {
		parent := app
		sitesMap := app.Sites

		for idx, name := range sites {
			var site *Site

			if _, p := sitesMap[name]; p == false {
				site = server.System.GetSite(*ctx, parent, name)
				site.Sites = make(map[string]*Site)
				site.Parent = parent
				site.server = server
				sitesMap[name] = site

				sitesMap = site.Sites
				parent = site
			} else {
				parent = site
				site = sitesMap[name]
				sitesMap = site.Sites
			}
			if len(sites) == idx+1 {
				ctx.Site = site
			}
		}
	}
}

func getModule(ctx *ServeContext, name string) {
	server := ctx.Server

	if name == "" {
		if ctx.GetConfig("modules.@0") != nil {
			name, _ = ctx.GetConfig("modules.@0").(string)
		}
	} else {
		if name != "_auth" {
			if ctx.GetConfig("modules") != nil {
				//check wheather part of config modules
				l, _ := ctx.GetConfig("modules.$length").(int)
				matched := false
				for i := 0; i < l; i++ {
					t, _ := ctx.GetConfig("modules.@" + strconv.Itoa(i)).(string)
					if t == name {
						matched = true
					}
				}
				if matched == false {
					name = ""
				}
			}
		}
	}

	if name != "" {
		if module, p := server.Modules[name]; p == false {
			ctx.Module = server.System.GetModule(*ctx, name)
			ctx.Module.mux = http.NewServeMux()
			ctx.Module.handlers = make(map[string]ServeHandler)
			ctx.Module.server = server
			ctx.Module.Build()
			server.Modules[name] = ctx.Module
		} else {
			ctx.Module = module
		}
	}
}

func (ctx *ServeContext) GetConfig(key string) interface{} {
	var value interface{}

	if ctx.Module != nil {
		value = ctx.Module.getConfig(key)
		if value != nil {
			return value
		}
	}

	site := ctx.Site

	for site != nil {
		value = site.getConfig(key)
		if value != nil {
			return value
		} else {
			if site.Parent != nil {
				site = site.Parent
			}
		}
	}

	if ctx.Application != nil {
		value = ctx.Application.getConfig(key)
		if value != nil {
			return value
		}
	}

	value = ctx.Server.getConfig(key)
	if value != nil {
		return value
	}
	return nil
}

func (ctx *ServeContext) GetJar() (*securecookie.SecureCookie, string) {
	var jar *securecookie.SecureCookie
	var prefix string
	if ctx.Application != nil && ctx.Application.jar != nil {
		jar = ctx.Application.jar
		prefix = ctx.Application.URI
	} else {
		jar = ctx.Server.jar
		prefix = "/"
	}

	return jar, prefix
}
