package httpserver

import (
	"golang-websocket-chat/internal/httpServer/middleware/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HttpServer struct {
	serveMux         *chi.Mux
	logger           *slog.Logger
	roomsHandlers    *roomsHandlers
	usersHandlers    *usersHandlers
	messagesHandlers *messagesHandlers
}

func New(rR roomsRepository, uR usersRepository, mR messagesRepository, log *slog.Logger) *HttpServer {
	router := chi.NewRouter()
	s := &HttpServer{
		logger:   log,
		serveMux: router,
		roomsHandlers: &roomsHandlers{
			repository: rR,
		},
		usersHandlers: &usersHandlers{
			repository: uR,
		},
		messagesHandlers: &messagesHandlers{
			repository: mR,
		},
	}
	s.serveMux.Use(middleware.RequestID)
	s.serveMux.Use(logger.New(s.logger))
	s.serveMux.Use(middleware.URLFormat)
	s.configureRoutes()
	return s
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *HttpServer) configureRoutes() {
	s.serveMux.Group(func(r chi.Router) {
		s.serveMux.Post("/room/subscribe", s.roomsHandlers.subscribeUserToRoom)
		s.serveMux.Post("/room/create", s.roomsHandlers.create)
	})
	s.serveMux.Group(func(r chi.Router) {
		s.serveMux.Post("/user/create", s.usersHandlers.create)
		s.serveMux.Delete("/user/delete/{id}", s.usersHandlers.delete)
		s.serveMux.Get("/user/{id}", s.usersHandlers.getById)
	})
	s.serveMux.Group(func(r chi.Router) {
		s.serveMux.Post("/messages/send", s.messagesHandlers.sendMessage)
		s.serveMux.Delete("/messages/delete/{id}", s.messagesHandlers.deleteMessage)
	})

}
