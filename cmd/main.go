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
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	r := gin.Default()

	r.POST("/check", api.Check)
	r.POST("/report", api.Report)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Println("Server started on port:", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Error graceful shutdown:", err)

		srv.Close()

	} else {
		config.Save()
		log.Println("Server stopped")
	}
}
