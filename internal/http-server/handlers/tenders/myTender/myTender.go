package my

import (
	"api/internal/models"
	"errors"
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
	Status      string `json:"status"`
	ServiceType string `json:"serviceType"`
	Version     uint64 `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

type ShowMyTender interface {
	GetMyTender(username string) (*models.Tender, error)
}

func New(log *slog.Logger, myTender ShowMyTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.my.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//limitStr := r.URL.Query().Get("limit")
		//offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		mytender, err := myTender.GetMyTender(username)
		if err != nil {
			log.Error("failed to get my tender", err.Error())
			render.JSON(w, r, errors.New("failed to get my tender"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Info("get my tender success")
		response := Response{
			Id:          mytender.Id.String(),
			Name:        mytender.Name,
			Description: mytender.Description,
			Status:      mytender.Status,
			ServiceType: mytender.ServiceType,
			Version:     mytender.Version,
			CreatedAt:   mytender.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
