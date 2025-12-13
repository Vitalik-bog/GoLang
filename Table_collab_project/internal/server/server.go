package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"table_collab/cmd/server/config"
	"table_collab/internal/server/ws"
	"table_collab/internal/service"
)

type Server struct {
	router *chi.Mux
	config *config.Config
	hub    *service.Hub
}

func New(cfg *config.Config) *Server {
	s := &Server{
		router: chi.NewRouter(),
		config: cfg,
		hub:    service.NewHub(cfg),
	}

	s.setupMiddleware()
	s.setupRoutes()

	go s.hub.Run()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	if s.config.Server.Env == "development" {
		s.router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}
}

func (s *Server) setupRoutes() {
	s.router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static"))))

	s.router.Get("/api/health", s.handleHealth)
	s.router.Get("/ws/{roomID}", s.handleWebSocket)
	s.router.Get("/", s.handleHome)
	s.router.Get("/room/{roomID}", s.handleRoomPage)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/index.html")
}

func (s *Server) handleRoomPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/room.html")
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	wsHandler := ws.NewHandler(s.hub)
	wsHandler.ServeWebSocket(roomID, w, r)
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         s.config.Server.Address,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.hub.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %v", err)
	}

	log.Println("Server stopped")
	return nil
}
