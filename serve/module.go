package serve

import (
	"net/http"

	"github.com/kirubasankars/metal"
)

type ServeHandler func(ctx ServeContext, w http.ResponseWriter, r *http.Request)

type Module struct {
	Name        string
	Path        string
	URI         string
	AuthEnabled bool
	AppEnabled  bool
	Config      *metal.Metal

	mux      *http.ServeMux
	handlers map[string]ServeHandler

	server *Server
}

func (module *Module) getConfig(key string) interface{} {
	if module.Config == nil {
		return nil
	}
	return module.Config.Get(key)
}

func (module *Module) Build() {
	if provider, p := moduleHanlderProvider[module.Name]; p {
		provider.Build(*module)
	} else {
		if provider, p := moduleHanlderProvider["."]; p {
			provider.Build(*module)
		}
	}
	mux := module.mux
	for pattern, handler := range module.handlers {
		if pattern != "" {
			var mh = new(ModuleHandler)
			mh.handler = handler
			mh.module = module
			uri := "/" + module.Name + pattern
			mux.Handle(uri, mh)

			if pattern == "/" || (pattern[len(pattern)-1:] == "/" && module.handlers[pattern[:len(pattern)-1]] == nil) {
				mux.HandleFunc(uri[:len(uri)-1], func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, r.URL.Path+"/", 301)
				})
			}
		}
	}
}

type ModuleHandler struct {
	handler ServeHandler
	module  *Module
}

func (mh *ModuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(mh.module.server, r.URL.Path)
	mh.handler(*ctx, w, r)
}
