package httpserver

import (
	"context"
	"encoding/json"
	"golang-websocket-chat/internal/models"
	resp "golang-websocket-chat/lib/api/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type usersRepository interface {
	Create(context.Context, *models.User) error
	Delete(context.Context, int) error
	GetById(context.Context, int) (*models.User, error)
}

type usersHandlers struct {
	repository usersRepository
}

func (h *usersHandlers) create(w http.ResponseWriter, r *http.Request) {
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
	usr := &models.User{
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
	type response struct {
		resp.Response
	}

	p := chi.URLParam(r, "id")
	id, err := strconv.Atoi(p)

	if err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}

	if err := h.repository.Delete(r.Context(), id); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}

	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}

func (h *usersHandlers) getById(w http.ResponseWriter, r *http.Request) {
	type response struct {
		resp.Response
		User *models.User
	}

	p := chi.URLParam(r, "id")
	id, err := strconv.Atoi(p)

	if err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}

	usr, err := h.repository.GetById(r.Context(), id)
	if err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
		User:     usr,
	})
}
