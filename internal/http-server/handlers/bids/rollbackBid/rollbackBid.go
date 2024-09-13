package rollbackBid

import (
	"api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

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
	reason string `json:"reason"`
}

type RollbackBid interface {
	RollbackBid(username string, version int, id string) (*models.Bid, error)
}

func New(log *slog.Logger, rollback RollbackBid) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.rollback.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "bidId")
		v := chi.URLParam(r, "version")
		version, err := strconv.Atoi(v)
		if err != nil {
			log.Error("failed to convert version", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		resp, err := rollback.RollbackBid(username, version, id)
		if err != nil {
			log.Error("failed to rollback bid", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("rollback bid success")
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
