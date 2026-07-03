package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/server"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SetUpRoutes()

	if err := database.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")); err != nil {
		logrus.Panicf("failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")

	go func() {
		if serverErr := srv.Run(":8080"); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", serverErr)

		}
	}()
	logrus.Print("server started at :8080")

	<-done

	logrus.Info("shutting down server")

	if dbCloseErr := database.ShutdownDatabase(); dbCloseErr != nil {
		logrus.WithError(dbCloseErr).Error("failed to close database connection")
	}

	if serverCloseErr := srv.Shutdown(shutDownTimeOut); serverCloseErr != nil {
		logrus.WithError(serverCloseErr).Panic("failed to gracefully shutdown server")
	}

}
