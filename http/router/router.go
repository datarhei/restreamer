package router

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Router interface {
	// FileRoutes
	FileRoutes() map[string]string

	// DirRoutes
	DirRoutes() map[string]string

	// StaticRoute
	StaticRoute() (string, string)
}

type router struct {
	prefixes    []string
	fileRoutes  map[string]string
	dirRoutes   map[string]string
	staticRoute string
}

func New(blockedPrefixes []string, routes map[string]string, static string) (Router, error) {
	r := &router{}

	r.prefixes = append(r.prefixes, blockedPrefixes...)
	r.staticRoute = static

	if len(r.staticRoute) != 0 {
		if path, err := filepath.Abs(r.staticRoute); err == nil {
			r.staticRoute = path
		} else {
			return nil, fmt.Errorf("invalid static route: %w", err)
		}

		r.prefixes = append(r.prefixes, "/ui")
	}

	r.fileRoutes = make(map[string]string)
	r.dirRoutes = make(map[string]string)

	if routes == nil {
		return r, nil
	}

	for route, target := range routes {
		for _, prefix := range r.prefixes {
			if strings.HasPrefix(route, prefix) {
				return nil, fmt.Errorf("the prefix of the route %s is blocked", route)
			}
		}

		if strings.HasSuffix(route, "/*") {
			route = strings.TrimSuffix(route, "/*")
			path, err := filepath.Abs(target)
			if err != nil {
				return nil, fmt.Errorf("invalid route %s: %w", target, err)
			}

			r.dirRoutes[route] = path
		} else {
			target = filepath.Clean(target)
			if !strings.HasPrefix(target, "/") {
				continue
			}

			r.fileRoutes[route] = target
		}
	}

	return r, nil
}

func (r *router) FileRoutes() map[string]string {
	routes := map[string]string{}

	for k, v := range r.fileRoutes {
		routes[k] = v
	}

	return routes
}

func (r *router) DirRoutes() map[string]string {
	routes := map[string]string{}

	for k, v := range r.dirRoutes {
		routes[k] = v
	}

	return routes
}

func (r *router) StaticRoute() (string, string) {
	return "/ui", r.staticRoute
}

func NewDummyRouter() Router {
	router, _ := New(nil, nil, "")

	return router
}

var DefaultRouter Router = NewDummyRouter()
