package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BI-Art-IT/ai-sdlc-backend/internal/config"
	"github.com/BI-Art-IT/ai-sdlc-backend/internal/database"
	"github.com/BI-Art-IT/ai-sdlc-backend/internal/migrate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Determine the path to migrations relative to the binary location.
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "migrations")

	// Run database migrations.
	if err = migrate.Up(cfg.DB.DSN(), migrationsPath); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	ctx := context.Background()

	pool, err := database.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("server listening on %s", addr)
	if err = http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
