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

type usersRepository interface {
	Create(context.Context, *postgres.User) error
	Delete(context.Context, int) error
	GetById(context.Context, int) (*postgres.User, error)
}

type usersHandlers struct {
	logger     *slog.Logger
	repository usersRepository
}

func (h *usersHandlers) create(w http.ResponseWriter, r *http.Request) {
	const op = "htmlServer.usersHandlers.create"

	type request struct {
		Name string `json:"name"`
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
	usr := &postgres.User{
		Name: req.Name,
	}
	if err := h.repository.Create(r.Context(), usr); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
		Name:     usr.Name,
	})
}

func (h *usersHandlers) delete(w http.ResponseWriter, r *http.Request) {
	const op = "htmlServer.usersHandlers.delete"
	type request struct {
		Id int `json:"id"`
	}
	type response struct {
		resp.Response
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	if err := h.repository.Delete(r.Context(), req.Id); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}

func (h *usersHandlers) getById(w http.ResponseWriter, r *http.Request) {
	const op = "htmlServer.usersHandlers.getById"
	type request struct {
		Id int `json:"id"`
	}
	type response struct {
		resp.Response
		User *postgres.User
	}
	req := &request{}
	usr := &postgres.User{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	usr, err := h.repository.GetById(r.Context(), req.Id)
	if err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
		User:     usr,
	})
}
