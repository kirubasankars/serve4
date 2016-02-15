package serve

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type CommonSiteBuilder struct{}

func (csb *CommonSiteBuilder) Build(module Module) {

	module.handlers["/"] = func(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
		fileName := ""
		url := r.URL.Path
		if ctx.Module.AppEnabled {
			fileName = filepath.Join(ctx.Module.Path, "app", url[ctx.URILength:])
		} else {
			fileName = filepath.Join(ctx.Module.Path, url[ctx.URILength:])
		}
		http.ServeFile(w, r, fileName)
	}

	module.handlers["/api/"] = func(ctx ServeContext, w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(ctx.Module.Path, r.URL.Path[ctx.URILength:], strings.ToLower(r.Method)+".json")
		f, _ := ioutil.ReadFile(path)
		w.Write(f)
	}
}
