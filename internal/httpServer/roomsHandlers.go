package httpserver

import (
	"context"
	"encoding/json"
	"golang-websocket-chat/internal/models"
	resp "golang-websocket-chat/lib/api/response"
	"net/http"

	"github.com/go-chi/render"
)

type roomsRepository interface {
	Create(context.Context, *models.Room) error
	Delete(context.Context, int) error
	GetById(context.Context, int) (*models.Room, error)
	AddUser(context.Context, int, int) error
}
type roomsHandlers struct {
	repository roomsRepository
}

func (h *roomsHandlers) subscribeUserToRoom(w http.ResponseWriter, r *http.Request) {
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
	if err := h.repository.AddUser(
		r.Context(),
		req.RoomId,
		req.UserId,
	); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}

func (h *roomsHandlers) create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name string `json:"room-name"`
	}
	type response struct {
		resp.Response
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	room := &models.Room{
		Name: req.Name,
	}
	if err := h.repository.Create(r.Context(), room); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}
