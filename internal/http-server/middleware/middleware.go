package apiMiddleware

import (
	"api/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type errResponse struct {
	Reason string `json:"reason"`
}

func CheckUserTenderMiddleware(log *slog.Logger, storage *postgresql.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.UserCheckMiddleware"
			log = log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			username := r.URL.Query().Get("username")
			if username == "" {
				log.Error("Username is empty")
				http.Error(w, "username is required", http.StatusUnauthorized)
			}
			id := chi.URLParam(r, "tenderId")

			var isAttached bool

			isAttached, err := storage.IsUserAttachedToOrganization(username, id)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Error("failed to found user", err.Error())
				render.JSON(w, r, errResponse{Reason: err.Error()})
				return
			}

			if !isAttached {
				w.WriteHeader(http.StatusForbidden)
				log.Error("user is not have permission", err.Error())
				render.JSON(w, r, errResponse{Reason: err.Error()})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
