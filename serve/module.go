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

func (module *Module) getConfig(key string) *string {
	if module.Config == nil {
		return nil
	}
	if v, e := module.Config.Get(key).(string); e == true {
		return &v
	} else {
		return nil
	}
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
			mux.HandleFunc("/"+module.Name+pattern, func(w http.ResponseWriter, r *http.Request) {
				ctx := NewContext(module.server, r.URL.Path)
				handler(*ctx, w, r)
			})

			if pattern == "/" {
				mux.HandleFunc("/"+module.Name, func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, r.URL.Path+"/", 301)
				})
			}
		}
	}
}
