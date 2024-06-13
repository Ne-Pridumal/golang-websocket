package httpserver

import (
	"context"
	"encoding/json"
	"golang-websocket-chat/internal/storage/postgres"
	resp "golang-websocket-chat/lib/api/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type roomsRepository interface {
	Create(context.Context, *postgres.Room) error
	Delete(context.Context, int) error
	GetById(context.Context, int) (*postgres.Room, error)
	AddUser(context.Context, int, int) error
}
type roomsHandlers struct {
	logger     *slog.Logger
	repository roomsRepository
}

func (h *roomsHandlers) subscribeUserToRoom(w http.ResponseWriter, r *http.Request) {
	const op = "htmlServer.roomsHandlers.subscribeUserToRoom"

	type request struct {
		UserId int `json:"user-id"`
		RoomId int `json:"room-id"`
	}
	type response struct {
		resp.Response
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	h.repository.AddUser(r.Context(), req.RoomId, req.UserId)
	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}

func (h *roomsHandlers) create(w http.ResponseWriter, r *http.Request) {
	const op = "htmlServer.roomsHandlers.createRoom"
	type request struct {
		Id   int    `json:"room-id"`
		Name string `json:"room-name"`
	}
	type response struct {
		resp.Response
		Name string
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	room := &postgres.Room{
		ID:   req.Id,
		Name: req.Name,
	}
	if err := h.repository.Create(r.Context(), room); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
		Name:     room.Name,
	})
}
