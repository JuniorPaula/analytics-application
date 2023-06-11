package usecases

import (
	"c2d-reports/internal/services"
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
	// calls chat2desk api
	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	for _, d := range dialogs {
		messages := chat2deskService.GetMessageByDialogID(d.ID)

		message := find(messages, func(msg services.Message) bool {
			return msg.Type == "to_client" || msg.Type == "from_client"
		})

		if message == nil {
			continue
		}

		var timerTMR int

		if message.Type == "from_client" {
			tmr, err := getTMR(message.Created)
			if err != nil {
				fmt.Println(err)
				continue
			}
			timerTMR = tmr
		} else {
			timerTMR = 0
		}

		client := chat2deskService.GetClientByID(message.ClientID)
		var statusTag string
		if len(client.Tags) > 0 {
			statusTag = client.Tags[0].Label
		} else {
			statusTag = "Sem tag"
		}

		for _, o := range operators {
			if o.OperatorID == message.OperatorID {
				fmt.Println("dialog id: ", d.ID)
				fmt.Println("operator id: ", o.OperatorID)
				fmt.Println("operator: ", o.OperatorName)
				fmt.Println("tmr: ", timerTMR)
				fmt.Println("qtd dialogs: ", o.QtdDialogs)
				fmt.Println("client phone: ", client.Phone)
				fmt.Println("status: ", statusTag)
				fmt.Println("------------------")
			}
		}

	}
}

func getTMR(created string) (int, error) {
	createdSplitTime := created[0:10]
	createParseTime, err := time.Parse("2006-01-02", createdSplitTime)
	if err != nil {
		return 0, err
	}

	todayTime := time.Now().UTC()

	timerTMR := todayTime.Sub(createParseTime).Seconds()

	return int(timerTMR), nil
}

func find(message []services.Message, predicate func(services.Message) bool) *services.Message {
	for _, m := range message {
		if predicate(m) {
			return &m
		}
	}
	return nil
}
