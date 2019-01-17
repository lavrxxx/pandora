package http

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"pandora/pkg/conf"
	"pandora/pkg/utils/errors"
	"pandora/pkg/utils/log"
)

// Route
type Route struct {
	Path       string
	Handler    func(w http.ResponseWriter, r *http.Request)
	Middleware []Middleware
	Method     string
}
type Routes []Route

// SubRoute
type SubRoute struct {
	Prefix     string
	Routes     Routes
	Middleware []Middleware
}
type SubRoutes []SubRoute

// Middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

const (
	Post = http.MethodPost
	Get  = http.MethodGet
)

// Listen start listen http requests
func Listen(endpoint string, subRoutes SubRoutes, static string) error {
	var r = mux.NewRouter()
	for _, subRoute := range subRoutes {
		s := r.PathPrefix(subRoute.Prefix).Subrouter()

		for _, route := range subRoute.Routes {
			middlewares := append(subRoute.Middleware, route.Middleware...)
			s.Handle(route.Path, handle(route.Handler, middlewares...)).Methods(route.Method)
		}
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	var h http.Handler
	h = handlers.LoggingHandler(os.Stdout, r)
	h = handlers.CORS(handlers.AllowedOrigins([]string{conf.Conf.Node.Endpoint}))(h)

	srv := &http.Server{
		Handler:           h,
		Addr:              endpoint,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 5,
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	log.Debugf("listen https server on %s", endpoint)
	return errors.WithStack(srv.ListenAndServeTLS(conf.Conf.TLS.Cert, conf.Conf.TLS.Key))
}

func handle(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	headers := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defaultHeaders(w, r)
			h.ServeHTTP(w, r)
		}
	}

	h = headers(h)

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func defaultHeaders(_ http.ResponseWriter, _ *http.Request) {}
