package http

import "net/http"

type Server struct {
	httpServer *http.Server
}

func New(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
