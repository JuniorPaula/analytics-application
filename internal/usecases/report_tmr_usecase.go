package usecases

import (
	"c2d-reports/internal/services"
	"encoding/json"
	"fmt"
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

	u.CalculateDialogsHanlder(dialogs)
}

func (u *ReportTmrUsecase) CalculateDialogsHanlder(dialogs []services.Dialog) {

	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	for _, d := range dialogs {
		requests := chat2deskService.GetDialogsByRequestID(d.LastMessage.RequestID)
		r, _ := json.Marshal(requests)
		fmt.Println(string(r))
	}
}

func getTMR(created string) int {
	createParseTime, _ := time.Parse("2006-01-02T15:04:05 MST", created)
	todayTime := time.Now()
	timerTMR := int(todayTime.Sub(createParseTime).Seconds())
	return timerTMR
}
