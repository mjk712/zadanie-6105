package myBids

import (
	"api/internal/models"
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

type ShowMyBids interface {
	GetMyBids(username string) ([]*models.Bid, error)
}

type errResponse struct {
	reason string `json:"reason"`
}

func New(log *slog.Logger, myBids ShowMyBids) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.myBids.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//limitStr := r.URL.Query().Get("limit")
		//offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		resp, err := myBids.GetMyBids(username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to get bids by username", err.Error())
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("get bid by username success")
		var response []Response
		for _, bid := range resp {
			response = append(response, Response{
				Id:         bid.Id.String(),
				Name:       bid.Name,
				Status:     bid.Status,
				AuthorType: bid.AuthorType,
				AuthorId:   bid.AuthorId.String(),
				Version:    bid.Version,
				CreatedAt:  bid.CreatedAt.Format(time.RFC3339),
			})
		}
		render.JSON(w, r, response)
	}
}
