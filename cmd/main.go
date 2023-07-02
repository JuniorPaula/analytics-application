package main

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/router"
	"c2d-reports/pkg/jobs"
	"c2d-reports/pkg/rabbitmq"
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
	// initialize .env variables
	config.InitVariables()

	// start the consumer on reports queue in a goroutine
	go rabbitmq.ConsumerOnReportsQueue()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// create the context to control of the exit
	_, cancel := context.WithCancel(context.Background())

	// create the server
	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: router.Handler()}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error while start the server: %v", err)
		}
	}()

	// print out the port
	fmt.Println("[::] Starting server on the port: ", config.Port)

	// schedule the jobs
	startSchedule()

	// wait for the interrupt signal
	<-interrupt
	log.Println("[::] Shutting down the server...")

	// cancel the context
	cancel()

	// wait for to shutdown the server
	timeout := 5 * time.Second
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), timeout)
	defer cancelShutdown()

	// shutdown the server
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("error while shutdown the server: %v", err)
	}

	log.Println("[::] Server gracefully stopped")
}

// startSchedule schedules the jobs
// 1. ScheduleCalculateReport - calculates the report for the current day
// 2. ScheduleDeleteReport - deletes the report for the current day
func startSchedule() {
	// schedule the jobs
	createReportSchedule := jobs.NewSchedule()
	go createReportSchedule.ScheduleCalculateReport()
	go createReportSchedule.ScheduleDeleteReport()
}
