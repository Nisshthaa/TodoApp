package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/server"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := server.SetUpRoutes()

	if err := database.OpenConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")

	go func() {
		if err := srv.Run(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", err)

		}
	}()
	logrus.Print("Server started at :8080")

	<-done

	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

}
