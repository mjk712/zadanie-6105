package tenderStatus

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type TenderStatus interface {
	GetTenderStatus(username string, id string) (string, error)
}
type errResponse struct {
	reason string `json:"reason"`
}

func New(log *slog.Logger, tendStat TenderStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.tenderStatus.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "tenderId")
		status, err := tendStat.GetTenderStatus(username, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to getTenders myTender tender", err.Error())
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("getTenders tender tenderStatus success")

		render.JSON(w, r, status)
	}
}
