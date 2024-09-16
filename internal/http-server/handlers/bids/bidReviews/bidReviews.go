package bidReviews

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	TenderId  string `json:"tenderId" db:"tender_id"`
	Feedback  string `json:"feedbackText" db:"feedback_text"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

type errResponse struct {
	Reason string `json:"reason"`
}

type ShowBidReviews interface {
	GetBidReviews(tenderId string, authorUsername string, requesterUsername string, limit int, offset int) ([]*Response, error)
}

func New(log *slog.Logger, bidReviews ShowBidReviews) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.bidReviews.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		tenderId := chi.URLParam(r, "tenderId")
		aName := r.URL.Query().Get("authorUsername")
		rName := r.URL.Query().Get("requesterUsername")

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

		resp, err := bidReviews.GetBidReviews(tenderId, aName, rName, limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to get bid reviews", err.Error())
			render.JSON(w, r, errResponse{Reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("get bid reviews success")

		render.JSON(w, r, resp)
	}
}
