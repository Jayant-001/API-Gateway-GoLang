package router

import (
	"jayant-001/api-gateway/internal/config"
	"jayant-001/api-gateway/internal/middleware"
	"jayant-001/api-gateway/internal/proxy"
	"net/http"
	"sync"
)

type Router struct {
	mu     sync.Mutex
	cfg    *config.Config
	router *http.ServeMux
}

func NewRouter(cfg *config.Config) (*Router, error) {
	r := &Router{
		cfg:    cfg,
		router: http.NewServeMux(), // a new HTTP request multiplexer
	}

	if err := r.updateRoutes(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) updateRoutes() error {

	r.mu.Lock()
	defer r.mu.Unlock()

	newRouter := http.NewServeMux()
	for _, route := range r.cfg.Routes {
		rp, err := proxy.NewReverseProxy(route.UpstreamURL)
		if err != nil {
			return err
		}

		var handler http.Handler = rp
		if route.AuthRequired {
			handler = middleware.Auth(r.cfg.Auth.APIKey)(handler)
		}
		newRouter.Handle(route.PathPrefix+"/", http.StripPrefix(route.PathPrefix, handler))
	}

	r.router = newRouter
	return nil
}

func (r *Router) ReloadConfig(cfg *config.Config) error {
	r.cfg = cfg
	return r.updateRoutes()
}
