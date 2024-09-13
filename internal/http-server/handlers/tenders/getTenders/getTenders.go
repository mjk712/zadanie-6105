package get

import (
	"api/internal/models"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

/*type Response struct {
	Status  string           `json:"status"`
	Error   string           `json:"error,omitempty"`
	Tenders []*models.Tender `json:"tenders"`
}*/

type TendersGetResponse interface {
	GetTenders() ([]*models.Tender, error)
}

func New(log *slog.Logger, response TendersGetResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tender.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		tenders, err := response.GetTenders()
		if err != nil {
			log.Error("failed to get tenders", err.Error())
			render.JSON(w, r, errors.New("failed to get tenders"))
			return
		}
		log.Info("get tenders success")
		render.JSON(w, r, tenders)
	}
}
