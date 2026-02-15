package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPublicRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", register)
		r.Post("/login", login)
		r.Post("/refresh", refresh)
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register ok"))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login ok"))
}

func refresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("refresh ok"))
}
