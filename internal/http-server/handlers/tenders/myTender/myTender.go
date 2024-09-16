package myTender

import (
	"api/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"tenderStatus"`
	ServiceType string `json:"serviceType"`
	Version     uint64 `json:"version"`
	CreatedAt   string `json:"createdAt"`
}
type errResponse struct {
	Reason string `json:"reason"`
}

type ShowMyTender interface {
	GetMyTender(username string, limit int, offset int) (*models.Tender, error)
}

func New(log *slog.Logger, myTender ShowMyTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.myTender.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")

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

		resp, err := myTender.GetMyTender(username, limit, offset)
		if err != nil {
			log.Error("failed to get my tender", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("get my tenders success")
		response := Response{
			Id:          resp.Id.String(),
			Name:        resp.Name,
			Description: resp.Description,
			Status:      resp.Status,
			ServiceType: resp.ServiceType,
			Version:     resp.Version,
			CreatedAt:   resp.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
