package add

import (
	"api/internal/models"
	"encoding/json"
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

type TenderNew interface {
	NewTender(tender *models.Tender) (*models.Tender, error)
}

func New(log *slog.Logger, newTender TenderNew) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.add.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var tender models.Tender
		err := json.NewDecoder(r.Body).Decode(&tender)
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			return
		}
		//todo проверка на существование пользователя
		resp, err := newTender.NewTender(&tender)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to create add tender", err.Error())
			render.JSON(w, r, errors.New("failed to create new tender"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Info("create add tender success")
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
