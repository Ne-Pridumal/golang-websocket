package websocket

import (
	"encoding/json"
	"golang-websocket-chat/internal/httpServer/middleware/logger"
	"golang-websocket-chat/internal/storage/postgres"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	storage  *postgres.Storage
	serveMux *chi.Mux
	logger   *slog.Logger
}

// func newChatServer() *chatServer {
// 	cs := &chatServer{
// 		subscriberMessageBuffer: 16,
// 		logf:                    log.Printf,
// 		subscribers:             make(map[*subscriber]struct{}),
// 		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
// 	}
// 	cs.serveMux.Handle("/", http.FileServer(http.Dir(".")))
// 	cs.serveMux.HandleFunc("/subscribe", cs.subscribeHandler)
// 	cs.serveMux.HandleFunc("/publish", cs.publishHandler)
//
// 	return cs
// }

func New(storage *postgres.Storage, log *slog.Logger) *Server {
	const op = "httpServer.ws.New"
	router := chi.NewRouter()
	s := &Server{
		storage:  storage,
		logger:   log,
		serveMux: router,
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
	s.serveMux.Post("/subscribe-room", s.subscribeUserToRoom)
	s.serveMux.Post("/create-room", s.createRoom)
}

func (s *Server) subscribeUserToRoom(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserId int `json:"user-id"`
		RoomId int `json:"room-id"`
	}
	const op = "htmlServer.ws.subscribeUserToRoom"
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.logger.Error("%s: %w", op, err)
		return
	}
	s.storage.Rooms().AddUser(r.Context(), req.RoomId, req.UserId)
}

func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID   int    `json:"room-id"`
		Name string `json:"room-name"`
	}
	const op = "htmlServer.ws.createRoom"
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.logger.Error("%s: %w", op, err)
		return
	}
	if err := s.storage.Rooms().Create(r.Context(), &postgres.Room{
		ID:   req.ID,
		Name: req.Name,
	}); err != nil {
		s.logger.Error("%s: %w", op, err)
		return
	}
}
