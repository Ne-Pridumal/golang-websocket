package httpserver

import (
	"golang-websocket-chat/internal/httpServer/middleware/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	serveMux      *chi.Mux
	logger        *slog.Logger
	roomsHandlers *roomsHandlers
	usersHandlers *usersHandlers
}

func New(rR roomsRepository, uR usersRepository, log *slog.Logger) *Server {
	const op = "httpServer.rs.New"
	router := chi.NewRouter()
	s := &Server{
		logger:   log,
		serveMux: router,
		roomsHandlers: &roomsHandlers{
			repository: rR,
			logger:     log,
		},
		usersHandlers: &usersHandlers{
			repository: uR,
			logger:     log,
		},
	}
	s.serveMux.Use(middleware.RequestID)
	s.serveMux.Use(logger.New(s.logger))
	s.serveMux.Use(middleware.URLFormat)
	s.configureRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) configureRoutes() {
	s.serveMux.Group(func(r chi.Router) {
		s.serveMux.Post("/room/subscribe", s.roomsHandlers.subscribeUserToRoom)
		s.serveMux.Post("/room/create", s.roomsHandlers.create)
	})
	s.serveMux.Group(func(r chi.Router) {
		s.serveMux.Post("/user/create", s.usersHandlers.create)
		s.serveMux.Post("/user/delete", s.usersHandlers.delete)
		s.serveMux.Post("/user/get-by-id", s.usersHandlers.getById)
	})

}
