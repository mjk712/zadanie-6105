package redactTender

import (
	"api/internal/models"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ServiceType *string `json:"serviceType"`
	Status      *string `json:"tenderStatus"`
}

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

type RedactTender interface {
	UpdateTender(username string, id string, updData *Request) (*models.Tender, error)
}

func New(log *slog.Logger, updTender RedactTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.updateTender.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "tenderId")

		var updData Request
		err := json.NewDecoder(r.Body).Decode(&updData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to decode request body", err.Error())
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		updatedTender, err := updTender.UpdateTender(username, id, &updData)
		if err != nil {
			log.Error("failed to update tender", err.Error())
			render.JSON(w, r, errors.New("failed to update tender"))
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("update tender success")
		response := Response{
			Id:          updatedTender.Id.String(),
			Name:        updatedTender.Name,
			Description: updatedTender.Description,
			Status:      updatedTender.Status,
			ServiceType: updatedTender.ServiceType,
			Version:     updatedTender.Version,
			CreatedAt:   updatedTender.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
