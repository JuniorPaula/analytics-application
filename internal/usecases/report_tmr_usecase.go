package usecases

import (
	"c2d-reports/internal/services"
	"fmt"
	"sync"
)

type ReportTmrUsecase struct {
	CompanyToken string
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

	operatorDialogData := make([]struct {
		OperatorID   int
		OperatorName string
		QtdDialogs   int
	}, len(operators))

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

	fmt.Println(operatorDialogData)
	fmt.Println("------------------")
	fmt.Println(dialogs)
}
