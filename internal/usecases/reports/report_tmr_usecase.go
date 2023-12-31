package usecases

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"c2d-reports/internal/services"
	"c2d-reports/pkg/rabbitmq"
	"fmt"
	"log"
	"sync"
	"time"
)

type ReportTmrUsecase struct {
	CompanyToken string
}

// LoadTMR loads the TMR report
func (u *ReportTmrUsecase) LoadTMR() {
	// difine wait group
	var wg sync.WaitGroup

	// calls chat2desk api
	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}
	operators := chat2deskService.GetOperators()

	var dialogs []services.Dialog
	dialogLock := sync.Mutex{}

	for i, o := range operators {
		wg.Add(1)
		go func(operatorID, index int) {
			defer wg.Done()
			d := chat2deskService.GetAllDialogsOpenByOperatorID(operatorID)

			// lock dialogs
			dialogLock.Lock()
			dialogs = append(dialogs, d...)
			dialogLock.Unlock()
		}(o.ID, i)
	}

	// wait for all goroutines to finish
	wg.Wait()

	u.CalculateDialogsHanlder(dialogs, operators)
}

// CalculateDialogsHanlder calculates the TMR for each dialog
// and sends the report to the queue
func (u *ReportTmrUsecase) CalculateDialogsHanlder(dialogs []services.Dialog, operators []services.Operator) {
	// define wait group
	var wg sync.WaitGroup

	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error while connect database;\n %s", err)
	}
	defer db.Close()

	var report repositories.Report

	for _, d := range dialogs {
		wg.Add(1)
		go func(dialog services.Dialog) {
			defer wg.Done()

			if dialog.End == "" {
				var tmrInSeconds int
				if dialog.LastMessage.Type == "from_client" {
					tmrInSeconds = getTMR(dialog.LastMessage.Created)
				} else {
					tmrInSeconds = 0
				}

				client := chat2deskService.GetClientByID(dialog.LastMessage.ClientID)
				statusTAG := "Sem tag"
				if len(client.Tags) > 0 {
					statusTAG = client.Tags[0].Label
				}

				for _, o := range operators {
					if o.ID == dialog.OperatorID {

						report.OperatorName = o.FirstName
						report.OperatorID = o.ID
						report.DialogID = dialog.ID
						report.TMRInSeconds = tmrInSeconds
						report.OpenedDialogs = o.OpenedDialogs
						report.Client = client.Phone
						report.StatusTAG = statusTAG

						rabbitmq.PusblisherOnReportsQueue(report)
					}
				}

			}

		}(d)

	}

	// wait for all goroutines to finish
	wg.Wait()
}

// getTMR returns time in seconds
// between created time and now
func getTMR(created string) int {
	createdParseTime, err := time.Parse("2006-01-02T15:04:05 MST", created)
	if err != nil {
		fmt.Println("Error while parse time")
		return 0
	}

	todayTime := time.Now()
	timerTMR := int(todayTime.Sub(createdParseTime).Seconds())
	return timerTMR
}
