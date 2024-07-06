package router

import (
	"net/http"
	"database/sql"
	s "restAPI/services"
	"context"
	"strings"

)


type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter(db *sql.DB) *Router {
	router := &Router{
		routes : make(map[string]map[string]http.HandlerFunc),
	}

	service := s.NewService(db)

    router.addRoute("POST", "/products", service.CreateProduct)
    router.addRoute("GET", "/products", service.ListProduct)
	router.addRoute("GET", "/products/{id}", service.GetProduct)

    return router
}


func (r *Router) addRoute(method string, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	for route, handler := range r.routes[method] {
		if match, params := matchRoute(route, path); match {
			ctx := req.Context()
			for k, v := range params {
				ctx = context.WithValue(ctx, k, v)
			}
			handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	http.NotFound(w, req)
}

func matchRoute(route, path string) (bool, map[string]string) {
	route_mtx := strings.Split(route, "/")
	path_mtx := strings.Split(path, "/")

	if len(route_mtx) != len(path_mtx) {
		return false, nil
	}

	params := make(map[string]string)
	for i := range route_mtx {
		if strings.HasPrefix(route_mtx[i], "{") && strings.HasSuffix(route_mtx[i], "}") {
			paramName := route_mtx[i][1 : len(route_mtx[i])-1]
			params[paramName] = path_mtx[i]
		} else if route_mtx[i] != path_mtx[i] {
			return false, nil
		}
	}
	return true, params
}