package changeTenderStatus

import (
	"api/internal/models"
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
	Status      string `json:"tenderStatus"`
	ServiceType string `json:"serviceType"`
	Version     uint64 `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

type errResponse struct {
	reason string `json:"reason"`
}

type ChangeTenderStatus interface {
	PutTenderStatus(username string, id string, status string) (*models.Tender, error)
}

func New(log *slog.Logger, changeStat ChangeTenderStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.changeTenderStatus.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		status := r.URL.Query().Get("tenderStatus")
		id := chi.URLParam(r, "tenderId")
		tend, err := changeStat.PutTenderStatus(username, id, status)
		if err != nil {
			log.Error("failed to change tender tenderStatus", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("change tender tenderStatus success")
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
