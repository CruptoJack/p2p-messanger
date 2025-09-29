package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/p2p-messanger/internal/handler"
	authmiddleware "github.com/p2p-messanger/internal/middleware"
	"github.com/p2p-messanger/internal/repository/postgresql"
	"github.com/p2p-messanger/internal/service"
	"github.com/p2p-messanger/pkg/config"
	"github.com/p2p-messanger/pkg/database"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("error load env")
	}

	jwtCfg := config.LoadJWT()
	if jwtCfg.SecretKey == "" {
		log.Fatal("JWT secret key is required")
	}

	dbCfg := database.ConfigFromENV()
	pool, err := database.NewPool(context.Background(), dbCfg)
	if err != nil {
		log.Fatal("Failed connect to database", err)
	}
	defer pool.Close()
	log.Println("Connect is successfully")
	userRepo := postgresql.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, jwtCfg)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	r.Route("/", func(r chi.Router) {
		authHandler.AuthRouters(r)
	})

	r.Route("/protected", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(authService))

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(authmiddleware.UserIDKey).(int64)
			response := map[string]interface{}{
				"message": "Access granted to protected resource",
				"user_id": userID,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server stating on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
