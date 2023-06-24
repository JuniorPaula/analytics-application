package main

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/router"
	"c2d-reports/pkg/jobs"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// initialize .env variables
	config.InitVariables()

	// print out the port
	fmt.Println("[::] Starting server on the port: ", config.Port)

	// initialize the router
	r := router.Handler()

	// schedule the jobs
	startSchedule()

	// start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))

}

func startSchedule() {
	// schedule the jobs
	createReportSchedule := jobs.NewSchedule()
	go createReportSchedule.ScheduleCalculateReport()
	go createReportSchedule.ScheduleDeleteReport()
}
