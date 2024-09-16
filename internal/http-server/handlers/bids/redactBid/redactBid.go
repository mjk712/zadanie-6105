package redactBid

import (
	"api/internal/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type Response struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"tenderStatus"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    uint64 `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type errResponse struct {
	Reason string `json:"reason"`
}

type RedactBid interface {
	UpdateBid(username string, id string, updData *Request) (*models.Bid, error)
}

func New(log *slog.Logger, updBid RedactBid) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.updateTender.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "bidId")

		var updData Request
		err := json.NewDecoder(r.Body).Decode(&updData)
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		resp, err := updBid.UpdateBid(username, id, &updData)
		if err != nil {
			log.Error("failed to update Bid", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("update bid success")
		response := Response{
			Id:         resp.Id.String(),
			Name:       resp.Name,
			Status:     resp.Status,
			AuthorType: resp.AuthorType,
			AuthorId:   resp.AuthorId.String(),
			Version:    resp.Version,
			CreatedAt:  resp.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
