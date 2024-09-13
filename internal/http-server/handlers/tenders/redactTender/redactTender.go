package redactTender

import (
	"api/internal/models"
	"errors"
	"github.com/go-chi/chi/v5"
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

type RedactTender interface {
	UpdateTender(username string, id string, status string) (*models.Tender, error)
}

func New(log *slog.Logger, changeTender RedactTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.redactTender.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "tenderId")
		tend, err := changeTender.RedactTender(username, id, status)
		if err != nil {
			log.Error("failed to change tender status", err.Error())
			render.JSON(w, r, errors.New("failed to change tender status"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Info("change tender status success")
		response := Response{
			Id:          tend.Id.String(),
			Name:        tend.Name,
			Description: tend.Description,
			Status:      tend.Status,
			ServiceType: tend.ServiceType,
			Version:     tend.Version,
			CreatedAt:   tend.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
