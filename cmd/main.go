package main

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/router"
	"c2d-reports/pkg/jobs"
	"c2d-reports/pkg/rabbitmq"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// initialize .env variables
	config.InitVariables()

	// start the consumer on reports queue in a goroutine
	go rabbitmq.ConsumerOnReportsQueue()

	// print out the port
	fmt.Println("[::] Starting server on the port: ", config.Port)

	// initialize the router
	r := router.Handler()

	// schedule the jobs
	startSchedule()

	// start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))

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
