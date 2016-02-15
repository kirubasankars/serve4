package serve

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/kirubasankars/metal"
)

type Server struct {
	port   string
	path   string
	Config *metal.Metal

	jar *securecookie.SecureCookie
	mux *http.ServeMux

	Apps    map[string]*Site
	Modules map[string]*Module

	System System
}

func (server *Server) Path() string {
	return server.path
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()

	fmt.Println("dsadasd", r.URL.Path)

	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	ctx := NewContext(server, r.URL.Path)

	mh := new(ServeHTTPHandler)
	mh.ServeHTTP(ctx, w, r)

	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

func (server *Server) getConfig(key string) interface{} {
	if server.Config == nil {
		return nil
	}
	return server.Config.Get(key)
}

func (server *Server) Start() {
	if err := http.ListenAndServe("localhost:"+server.port, server.mux); err != nil {
		fmt.Println(err)
	}
}

func NewServer(path string, port string) *Server {
	server := new(Server)

	server.path = path
	server.port = port
	server.Apps = make(map[string]*Site)
	server.Modules = make(map[string]*Module)
	server.Config = getConfig(server.path)

	server.mux = http.NewServeMux()
	server.mux.Handle("/", server)

	return server
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
