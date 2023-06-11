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

	operatorDialogData := make([]OperatorDialogData, len(operators))

	opDialogLock := sync.Mutex{}

	for i, o := range operators {
		wg.Add(2)
		go func(operatorID, index int) {
			defer wg.Done()
			d := chat2deskService.GetAllDialogsOpenByOperatorID(operatorID)

			// lock dialogs
			dialogLock.Lock()
			dialogs = append(dialogs, d...)
			dialogLock.Unlock()
		}(o.ID, i)

		go func(operatorID int, operatorName string, index int) {
			defer wg.Done()
			d := chat2deskService.GetDialogsByOperator(operatorID)

			// lock the access to operatorDialogData
			opDialogLock.Lock()
			operatorDialogData[index].OperatorID = operatorID
			operatorDialogData[index].OperatorName = operatorName
			operatorDialogData[index].QtdDialogs = len(d)
			opDialogLock.Unlock()
		}(o.ID, o.FirstName, i)
	}

	// wait for all goroutines to finish
	wg.Wait()

	u.MessageMapperHanlder(dialogs, operatorDialogData)
}

func (u *ReportTmrUsecase) MessageMapperHanlder(dialogs []services.Dialog, operators []OperatorDialogData) {
	// define wait group
	var wg sync.WaitGroup

	// calls chat2desk api
	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	var dialogInfo DialogInfo

	for _, d := range dialogs {
		wg.Add(1)
		go func(dialog services.Dialog) {
			defer wg.Done()
			messages := chat2deskService.GetMessageByDialogID(dialog.ID)

			message := findMessageByType(messages, "from_client", "from_operator")
			if message == nil {
				return
			}

			timerTMR := 0
			if message.Type == "from_client" {
				tmr := getTMR(message.Created)
				timerTMR = tmr
			}

			client := chat2deskService.GetClientByID(message.ClientID)
			statusTag := "Sem tag"
			if len(client.Tags) > 0 {
				statusTag = client.Tags[0].Label
			}

			for _, o := range operators {
				if o.OperatorID == message.OperatorID {
					dialogInfo = DialogInfo{
						DialogID:     dialog.ID,
						OperatorID:   o.OperatorID,
						OperatorName: o.OperatorName,
						TmrInSeconds: timerTMR,
						QtdDialogs:   o.QtdDialogs,
						ClientPhone:  client.Phone,
						Status:       statusTag,
					}

					jsonData, err := json.Marshal(dialogInfo)
					if err != nil {
						fmt.Println("Error while convert to json:", err)
						return
					}
					fmt.Println("json data: ", string(jsonData))
				}
			}
		}(d)

	}

	// wait for all goroutines to finish
	wg.Wait()

}

func getTMR(created string) int {
	createParseTime, _ := time.Parse("2006-01-02T15:04:05 MST", created)
	todayTime := time.Now()
	timerTMR := int(todayTime.Sub(createParseTime).Seconds())
	return timerTMR
}

func findMessageByType(messages []services.Message, types ...string) *services.Message {
	for _, message := range messages {
		for _, t := range types {
			if message.Type == t {
				return &message
			}
		}
	}
	return nil
}
