package addBid

import (
	"api/internal/models"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
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
	Reason string `json:"reason"`
}

type BidNew interface {
	NewBid(bid *models.Bid) (*models.Bid, error)
}

func New(log *slog.Logger, newBid BidNew) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.addBid.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var bid models.Bid
		err := json.NewDecoder(r.Body).Decode(&bid)
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			return
		}

		resp, err := newBid.NewBid(&bid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to add bid", err.Error())
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("add bid success")
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
