package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutesExist(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()

	chiRoutes, ok := testRoutes.(chi.Router)

	if !ok {
		t.Error("Error: unable to cast testRoutes to chiRoutes")
	}

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Error("Error: did not find " + route + " in registered routes.")
	}
}
