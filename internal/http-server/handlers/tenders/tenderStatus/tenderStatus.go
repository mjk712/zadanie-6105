package status

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type TenderStatus interface {
	GetTenderStatus(username string, id string) (string, error)
}

func New(log *slog.Logger, tendStat TenderStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.myTender.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//limitStr := r.URL.Query().Get("limit")
		//offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "tenderId")
		status, err := tendStat.GetTenderStatus(username, id)
		if err != nil {
			log.Error("failed to getTenders myTender tender", err.Error())
			render.JSON(w, r, errors.New("failed to getTenders myTender tender"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Info("getTenders tender status success")

		render.JSON(w, r, status)
	}
}
