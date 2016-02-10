package serve

import (
	"fmt"
	"net/http"
	"path/filepath"
)

var moduleHanlderProvider map[string]ModuleHandlerProvider

type ModuleHandlerProvider interface {
	Build(module Module)
}

func RegisterHandlerProvider(name string, mhr ModuleHandlerProvider) {
	moduleHanlderProvider[name] = mhr
}

func init() {
	moduleHanlderProvider = make(map[string]ModuleHandlerProvider)
	RegisterHandlerProvider(".", new(CommonSiteHandlerProvider))
}

type CommonSiteHandlerProvider struct{}

func (cshp *CommonSiteHandlerProvider) Build(module Module) {
	module.handlers["/"] = func(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
		fileName := ""
		url := r.URL.Path
		if ctx.Module.AppEnabled {
			fileName = filepath.Join(ctx.Module.Path, "app", url[ctx.URILength:])
		} else {
			fileName = filepath.Join(ctx.Module.Path, url[ctx.URILength:])
		}
		fmt.Println(fileName)
		http.ServeFile(w, r, fileName)
	}
}
