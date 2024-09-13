package ping

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ping"
		w.Header().Set("Content-Type", "application/json")
		pong := "Pong"
		err := json.NewEncoder(w).Encode(pong)
		if err != nil {
			log.Error(op, "failed to pong", err)
		}
		w.WriteHeader(http.StatusOK)
		log.Info("ping", slog.String("pong", pong))
	}
}
