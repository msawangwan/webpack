package webpack

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Exported errors
var (
	ErrCouldNotRegisterUnsupportedMethod = errors.New("handler registration failed: unsupported method")
)

type handlers struct {
	registered map[string]http.HandlerFunc
	patterns   map[string]*regexp.Regexp
}

func newHandlers() *handlers {
	return &handlers{
		registered: make(map[string]http.HandlerFunc),
		patterns:   make(map[string]*regexp.Regexp),
	}
}

// Router manages handlers and regex path resolution
type Router struct {
	Routes map[string]*handlers
}

// NewRouter returns a Router struct that by default supports four HTTP method verbs (GET, POST, PUT, DELETE)
func NewRouter() *Router {
	return &Router{
		Routes: map[string]*handlers{
			"get":    newHandlers(),
			"post":   newHandlers(),
			"put":    newHandlers(),
			"delete": newHandlers(),
		},
	}
}

// GET is a wrapper around Register that registers the pattern 'route' to the GET handler 'h'
func (router *Router) GET(route string, h http.HandlerFunc) error {
	return router.Register(http.MethodGet, route, h)
}

// POST is a wrapper around Register that registers the pattern 'route' to the POST handler 'h'
func (router *Router) POST(route string, h http.HandlerFunc) error {
	return router.Register(http.MethodGet, route, h)
}

// PUT is a wrapper around Register that registers the pattern 'route' to the PUT handler 'h'
func (router *Router) PUT(route string, h http.HandlerFunc) error {
	return router.Register(http.MethodGet, route, h)
}

// DELETE is a wrapper around Register that registers the pattern 'route' to the DELETE handler 'h'
func (router *Router) DELETE(route string, h http.HandlerFunc) error {
	return router.Register(http.MethodGet, route, h)
}

// Register registers a handler (or another router) to 'route', the specified url pattern
func (router *Router) Register(method, route string, h http.HandlerFunc) error {
	m := strings.ToLower(method)

	if _, e := router.Routes[m]; e {
		re, err := regexp.Compile(route)
		if err != nil {
			return err
		}

		router.Routes[m].registered[route] = h
		router.Routes[m].patterns[route] = re

		return nil
	}

	return ErrCouldNotRegisterUnsupportedMethod
}

// ServeHTTP searches registered handlers that match on an incoming url pattern
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path

	log.Printf("%s - %s", r.Method, u)

	m := strings.ToLower(r.Method)

	if routes, exists := router.Routes[m]; exists {
		for pattern, handler := range routes.registered {
			exists := routes.patterns[pattern].MatchString(u)

			if exists {
				handler(w, r)
				return
			}
		}
	}

	http.NotFound(w, r) // TODO: serve a default not found page
}

// EntrypointRouter is intended for serving entrypoints defined in a webpack config but can serve any static file
type EntrypointRouter struct {
	entrypoints map[string]string
}

// NewEntrypointRouter returns an initialized EntrypointRouter struct
func NewEntrypointRouter() *EntrypointRouter {
	return &EntrypointRouter{
		entrypoints: make(map[string]string),
	}
}

// RegisterFile registers the url 'route' to a static file at 'filepath'
func (router *EntrypointRouter) RegisterFile(route, filepath string) {
	router.entrypoints[route] = filepath
}

// ServeHTTP serves the entrypoint registered to a specific route
func (router *EntrypointRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("(entrypoint) %s - %s", r.Method, r.URL.Path)

	if r.Method == http.MethodGet {
		if f, e := router.entrypoints[r.URL.Path]; e {
			http.ServeFile(w, r, f)
		} else {
			http.NotFound(w, r)
		}
	} else {
		http.Error(w, fmt.Sprintf("cannot %s %s", r.Method, r.URL.Path), http.StatusMethodNotAllowed)
	}
}

// ResourceRouter is intended to serve webpack assets but is simply a fileserver
type ResourceRouter struct {
	resources map[string]http.Handler
}

// NewResourceRouter returns an initialised ResourceRouter struct
func NewResourceRouter() *ResourceRouter {
	return &ResourceRouter{
		resources: make(map[string]http.Handler),
	}
}

// RegisterDirectory adds a directory into a list of directories that contain public assets
func (router *ResourceRouter) RegisterDirectory(dir string) {
	router.resources[dir] = http.FileServer(http.Dir(dir))
}

// ServeHTTP checks each registered directory for the requested resource
func (router *ResourceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("(resource) %s - %s", r.Method, r.URL.Path)

	for root, handler := range router.resources {
		_, err := os.Stat(filepath.Join(root, r.URL.Path))
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, fmt.Sprintf("%s", err), 404)
				return
			}
		}

		handler.ServeHTTP(w, r)

		return
	}

	http.NotFound(w, r)
}
