package bidFeedback

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
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"tenderStatus"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    uint64 `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type WriteBidFeedback interface {
	PutBidFeedback(username string, id string, feedback string) (*models.Bid, error)
}

type errResponse struct {
	reason string `json:"reason"`
}

func New(log *slog.Logger, putFeedback WriteBidFeedback) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.bids.writeBidFeedback.New"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		username := r.URL.Query().Get("username")
		feedback := r.URL.Query().Get("feedback")
		id := chi.URLParam(r, "bidId")
		resp, err := putFeedback.PutBidFeedback(username, id, feedback)
		if err != nil {
			log.Error("failed to put bid feedback", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse{reason: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("put bid feedback success")
		response := Response{
			Id:         resp.Id.String(),
			Name:       resp.Name,
			Status:     resp.Status,
			AuthorType: resp.AuthorType,
			AuthorId:   resp.AuthorId.String(),
			Version:    resp.Version,
			CreatedAt:  resp.CreatedAt.Format(time.RFC3339),
		}

		render.JSON(w, r, response)
	}
}
