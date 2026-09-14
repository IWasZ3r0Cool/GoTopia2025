package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IWasZ3r0Cool/GoTopia2025/internal/game"
	gameserver "github.com/IWasZ3r0Cool/GoTopia2025/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	engine := game.NewGameEngine()
	hub := gameserver.NewHub(engine)
	go hub.Run(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	staticDir := os.Getenv("GOTOPIA_STATIC_DIR")
	if staticDir == "" {
		staticDir = "frontend/dist"
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           gameserver.NewHTTPHandler(hub, staticDir),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown server: %v", err)
		}
	}()

	log.Printf("GoTopia server listening on http://localhost:%s", port)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
