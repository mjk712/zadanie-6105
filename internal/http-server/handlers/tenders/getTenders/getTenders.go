package getTenders

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
	reason string `json:"reason"`
}

type TendersGetResponse interface {
	GetTenders(limit int, offset int, serviceTypes []string) ([]*models.Tender, error)
}

func New(log *slog.Logger, response TendersGetResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.getTenders.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		queryParams := r.URL.Query()

		limitStr := queryParams.Get("limit")
		offsetStr := queryParams.Get("offset")
		serviceTypes := queryParams["serviceType"]

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

		resp, err := response.GetTenders(limit, offset, serviceTypes)
		if err != nil {
			log.Error("failed to get tenders", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}
		log.Info("get tenders success")
		w.WriteHeader(http.StatusOK)

		var response []Response
		for _, tender := range resp {
			response = append(response, Response{
				Id:          tender.Id.String(),
				Name:        tender.Name,
				Description: tender.Description,
				Status:      tender.Status,
				ServiceType: tender.ServiceType,
				Version:     tender.Version,
				CreatedAt:   tender.CreatedAt.Format(time.RFC3339),
			})
		}
		render.JSON(w, r, response)
	}
}
