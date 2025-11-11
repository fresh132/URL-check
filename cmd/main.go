package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fresh132/URL-check/internal/api"
	"github.com/fresh132/URL-check/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	config.Load()

	r := gin.Default()

	r.POST("/check", api.Check)
	r.POST("/report", api.Report)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Println("Server started on port:", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Error graceful shutdown:", err)

		srv.Close()

	} else {
		config.Save()
		log.Println("Server stopped")
	}
}
