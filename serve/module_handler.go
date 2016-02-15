package serve

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

type ServeHTTPHandler struct{}

func (serveHTTPHandler *ServeHTTPHandler) ServeHTTP(ctx *ServeContext, w http.ResponseWriter, r *http.Request) {

	if ctx == nil || ctx.Module == nil {
		http.NotFound(w, r)
		return
	}

	var prefix string
	var jar *securecookie.SecureCookie

	if ctx.Application != nil {
		prefix = ctx.Application.URI
		jar = ctx.Application.jar
	} else {
		prefix = ""
		jar = ctx.Server.jar
	}

	if ctx.Module.AuthEnabled {
		if cookie, err := r.Cookie("_auth"); err == nil {
			value := make(map[string]string)
			if err = jar.Decode("_auth", cookie.Value, &value); err == nil {
				serveHTTPHandler.Serve(ctx, w, r)
			} else {
				http.Redirect(w, r, prefix+"/_auth?redirectUrl="+r.URL.Path, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, prefix+"/_auth?redirectUrl="+r.URL.Path, http.StatusFound)
		}
	} else {
		serveHTTPHandler.Serve(ctx, w, r)
	}
}

func (serveHTTPHandler *ServeHTTPHandler) Serve(ctx *ServeContext, w http.ResponseWriter, r *http.Request) {
	newR, _ := http.NewRequest(r.Method, "/"+ctx.Module.Name+r.URL.Path[ctx.URILength:], nil)
	handler, _ := ctx.Module.mux.Handler(newR)
	newR.URL.Path = r.URL.Path
	handler.ServeHTTP(w, r)
}
