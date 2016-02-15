package serve

import (
	"github.com/gorilla/securecookie"
	"github.com/kirubasankars/metal"
)

type Site struct {
	Name   string
	Path   string
	URI    string
	Config *metal.Metal

	Sites map[string]*Site
	jar   *securecookie.SecureCookie

	Parent *Site
	server *Server
}

func (site *Site) getConfig(key string) interface{} {
	if site.Config == nil {
		return nil
	}
	return site.Config.Get(key)
}
