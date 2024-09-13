package rollbackTender

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

type RollbackTender interface {
	RollbackTender(username string, version int, id string) (*models.Tender, error)
}

func New(log *slog.Logger, rollback RollbackTender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.rollback.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "tenderId")
		v := chi.URLParam(r, "version")
		version, err := strconv.Atoi(v)
		if err != nil {
			log.Error("failed to convert version", err.Error())
			return
		}

		rollbackedTender, err := rollback.RollbackTender(username, version, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to rollback tender", err.Error())
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("rollback tender success")
		response := Response{
			Id:          rollbackedTender.Id.String(),
			Name:        rollbackedTender.Name,
			Description: rollbackedTender.Description,
			Status:      rollbackedTender.Status,
			ServiceType: rollbackedTender.ServiceType,
			Version:     rollbackedTender.Version,
			CreatedAt:   rollbackedTender.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
