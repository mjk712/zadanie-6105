package bidStatus

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type BidStatus interface {
	GetBidStatus(username string, id string) (string, error)
}

type errResponse struct {
	Reason string `json:"reason"`
}

func New(log *slog.Logger, bidStat BidStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		const op = "handlers.bids.bidStatus.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//limitStr := r.URL.Query().Get("limit")
		//offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		id := chi.URLParam(r, "bidId")
		status, err := bidStat.GetBidStatus(username, id)
		if err != nil {
			log.Error("failed to get bid status", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("get bid status success")

		render.JSON(w, r, status)
	}
}
