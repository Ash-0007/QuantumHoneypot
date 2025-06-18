package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"pqcd/api"
	"pqcd/security"
)

func main() {
	// Parse command line flags
	var (
		port        = flag.Int("port", 8080, "Port to listen on")
		enableAI    = flag.Bool("enable-ai", false, "Enable AI threat detection")
		logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Configure logging
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Create router
	r := mux.NewRouter()

	// Initialize API routes
	api.RegisterRoutes(r)
	
	// Initialize AI security if enabled
	if *enableAI {
		logrus.Info("Initializing AI security layer")
		aiHandler := security.NewAISecurityMiddleware()
		r.Use(aiHandler.Middleware)
	}

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.ExposedHeaders([]string{"X-Anomaly-Detected", "X-Anomaly-Score"}),
	)
	
	// Configure server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler(r),
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on port %d", *port)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	logrus.Info("Server shutdown complete")
} 