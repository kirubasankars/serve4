package serve

import "net/http"

type ServeHTTPHandler struct{}

func (serveHTTPHandler *ServeHTTPHandler) ServeHTTP(ctx *ServeContext, w http.ResponseWriter, r *http.Request) {

	fa := new(FormsAuthentication)

	if ctx == nil || ctx.Module == nil {
		http.NotFound(w, r)
		return
	}

	if ctx.Module.AuthEnabled {
		if fa.Validate(ctx, w, r) {
			serveHTTPHandler.Serve(ctx, w, r)
		} else {
			fa.RedirectToLogin(*ctx, w, r)
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
