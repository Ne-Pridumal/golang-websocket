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

type messagesRepository interface {
	Create(context.Context, *models.Message) error
	Delete(context.Context, int) error
}

type messagesHandlers struct {
	repository messagesRepository
}

func (h *messagesHandlers) sendMessage(w http.ResponseWriter, r *http.Request) {
	type request struct {
		RoomId  int    `json:"room-id"`
		Content string `json:"content"`
	}

	type response struct {
		resp.Response
	}

	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	if err := h.repository.Create(
		r.Context(),
		&models.Message{
			RoomId:  req.RoomId,
			Content: req.Content,
		},
	); err != nil {
		render.JSON(w, r, resp.Error(err.Error()))
		return
	}
	render.JSON(w, r, response{
		Response: resp.OK(),
	})
}

func (h *messagesHandlers) deleteMessage(w http.ResponseWriter, r *http.Request) {
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
