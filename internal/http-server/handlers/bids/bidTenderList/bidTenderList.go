package bidTenderList

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

type ShowBidTenderList interface {
	GetBidTenderList(tenderId string, limit int, offset int) ([]*models.Bid, error)
}

type errResponse struct {
	Reason string `json:"reason"`
}

func New(log *slog.Logger, bidTenderList ShowBidTenderList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.bidTenderList.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		tenderId := chi.URLParam(r, "tenderId")

		limit := 10
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			}
		}

		resp, err := bidTenderList.GetBidTenderList(tenderId, limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to get bids by tenderId", err.Error())
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("get bids by tenderId success")
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
