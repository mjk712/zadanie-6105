package addTender

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

type TenderNew interface {
	NewTender(tender *models.Tender) (*models.Tender, error)
}

func New(log *slog.Logger, newTender TenderNew) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.addTender.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var tender models.Tender
		err := json.NewDecoder(r.Body).Decode(&tender)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to decode request body", err.Error())
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		resp, err := newTender.NewTender(&tender)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to create addTender tender", err.Error())
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Info("create new tender success")
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
