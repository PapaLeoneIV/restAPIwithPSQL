package router

import (
	"net/http"
	"database/sql"
	s "restAPI/services"
	"context"
	"strings"
	"fmt"
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
    router.addRoute("GET", "/products", service.GetProducts)
	router.addRoute("GET", "/products/{id}", service.GetProducts)

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
				fmt.Printf("ctx %v\n", ctx)	
			}
			handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	http.NotFound(w, req)
}

func matchRoute(route, path string) (bool, map[string]string) {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i := range routeParts {
		if strings.HasPrefix(routeParts[i], "{") && strings.HasSuffix(routeParts[i], "}") {
			paramName := routeParts[i][1 : len(routeParts[i])-1]
			params[paramName] = pathParts[i]
		} else if routeParts[i] != pathParts[i] {
			return false, nil
		}
	}
	fmt.Printf("params : %v\n", params)
	return true, params
}