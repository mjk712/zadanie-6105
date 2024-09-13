package main

import (
	"api/internal/config"
	"api/internal/http-server/handlers/bids/addBid"
	"api/internal/http-server/handlers/bids/bidFeedback"
	"api/internal/http-server/handlers/bids/bidReviews"
	"api/internal/http-server/handlers/bids/bidStatus"
	"api/internal/http-server/handlers/bids/bidTenderList"
	"api/internal/http-server/handlers/bids/changeBidStatus"
	"api/internal/http-server/handlers/bids/myBids"
	"api/internal/http-server/handlers/bids/redactBid"
	"api/internal/http-server/handlers/bids/rollbackBid"
	"api/internal/http-server/handlers/bids/submitBidDecision"
	"api/internal/http-server/handlers/ping"
	"api/internal/http-server/handlers/tenders/addTender"
	"api/internal/http-server/handlers/tenders/changeTenderStatus"
	"api/internal/http-server/handlers/tenders/getTenders"
	"api/internal/http-server/handlers/tenders/myTender"
	"api/internal/http-server/handlers/tenders/redactTender"
	"api/internal/http-server/handlers/tenders/rollbackTender"
	"api/internal/http-server/handlers/tenders/tenderStatus"
	"api/internal/storage/postgresql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	//config
	cfg := config.New()
	fmt.Println(cfg)
	//log
	log := setupLogger(cfg.Env)
	log.Info(
		"starting tender api",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)
	log.Debug("debug messages are enabled")
	//storage
	storage, err := postgresql.New(cfg.ConnectionString)
	if err != nil {
		log.Error("failed to init storage", err)
	}
	//router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Get("/ping", ping.New(log))
		r.Get("/tenders", getTenders.New(log, storage))
		r.Post("/tenders/new", addTender.New(log, storage))
		r.Get("/tenders/my", myTender.New(log, storage))
		r.Get("/tenders/{tenderId:[a-fA-f0-9\\-]{36}}/status", tenderStatus.New(log, storage))
		r.Put("/tenders/{tenderId:[a-fA-f0-9\\-]{36}}/status", changeTenderStatus.New(log, storage))
		r.Patch("/tenders/{tenderId:[a-fA-f0-9\\-]{36}}/edit", redactTender.New(log, storage))
		r.Put("/tenders/{tenderId:[a-fA-f0-9\\-]{36}}/rollback/{version}", rollbackTender.New(log, storage))

		r.Post("/bids/new", addBid.New(log, storage))
		r.Get("/bids/myBids", myBids.New(log, storage))
		r.Get("/bids/{tenderId}/list", bidTenderList.New(log, storage))
		r.Get("/bids/{bidId}/status", bidStatus.New(log, storage))
		r.Put("/bids/{bidId}/status", changeBidStatus.New(log, storage))
		r.Patch("/bids/{bidId}/edit", redactBid.New(log, storage))
		r.Put("/bids/{bidId}/submit_decision", submitBidDecision.New(log, storage))
		r.Put("/bids/{bidId}/feedback", bidFeedback.New(log, storage))
		r.Put("/bids/{bidId}/rollback/{version}", rollbackBid.New(log, storage))
		r.Get("/bids/{tenderId}/reviews", bidReviews.New(log, storage))

	})
	http.ListenAndServe(":8080", router)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
