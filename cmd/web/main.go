package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	application := InitializedApp()
	app := application.App
	routeConfig := application.RouteConfig

	routeConfig.Setup()

	host := app.Config.GetString("APP_HOST")
	port := app.Config.GetInt("APP_PORT")
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:    addr,
		Handler: app.Gin,
	}

	go func() {
		log.Printf("Server starting on: %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")

	db, err := app.Database.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
	log.Println("Database connection closed")
}
