package usecases

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"c2d-reports/internal/services"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type ReportTmrUsecase struct {
	CompanyToken string
}

type OperatorDialogData struct {
	OperatorID   int
	OperatorName string
	QtdDialogs   int
}

type DialogInfo struct {
	DialogID     int    `json:"dialog_id"`
	OperatorID   int    `json:"operator_id"`
	OperatorName string `json:"operator_name"`
	TmrInSeconds int    `json:"tmr_in_seconds"`
	QtdDialogs   int    `json:"qtd_dialogs"`
	ClientPhone  string `json:"client_phone"`
	Status       string `json:"status"`
}

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

func (u *ReportTmrUsecase) CalculateDialogsHanlder(dialogs []services.Dialog, operators []services.Operator) {
	// define wait group
	var wg sync.WaitGroup

	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error while connect database; %s", err)
	}
	defer db.Close()

	repository := repositories.NewReportRepository(db)

	var report repositories.Report

	for _, d := range dialogs {
		wg.Add(1)
		go func(dialog services.Dialog) {
			defer wg.Done()

			requests := chat2deskService.GetDialogsByRequestID(dialog.LastMessage.RequestID)

			sort.Sort(sortByCreated(requests))

			messageFromClient := findMessageIN(requests)

			var tmrInSeconds int
			if messageFromClient != nil {
				tmrInSeconds = getTMR(messageFromClient.Created)
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
					report.Status = statusTAG

					reportID, err := repository.CreateOrUpdate(report)
					if err != nil {
						log.Fatalf("Error while create or update report; %s", err)
					}
					fmt.Printf("ID: [%d]; report computed:", reportID)
					fmt.Println("--------------------------")
				}
			}

		}(d)

	}

	// wait for all goroutines to finish
	wg.Wait()
}

func findMessageIN(requests []services.ResponseRequests) *services.ResponseRequests {
	for _, r := range requests {
		if r.Type == "in" {
			return &requests[0]
		}
	}
	return nil
}

type sortByCreated []services.ResponseRequests

func (r sortByCreated) Len() int           { return len(r) }
func (r sortByCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }
func (r sortByCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func getTMR(created int) int {
	createdParseTime := time.Unix(int64(created), 0)
	todayTime := time.Now()
	timerTMR := int(todayTime.Sub(createdParseTime).Seconds())
	return timerTMR
}
