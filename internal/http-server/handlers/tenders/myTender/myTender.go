package myTender

import (
	"api/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
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
	reason string `json:"reason"`
}

type ShowMyTender interface {
	GetMyTender(username string) (*models.Tender, error)
}

func New(log *slog.Logger, myTender ShowMyTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.myTender.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//limitStr := r.URL.Query().Get("limit")
		//offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		resp, err := myTender.GetMyTender(username)
		if err != nil {
			log.Error("failed to getTenders myTender tender", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("getTenders myTender tender success")
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
