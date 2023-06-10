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

	for _, o := range operators {
		wg.Add(1)
		go func(operatorID int) {
			defer wg.Done()
			d := chat2deskService.GetAllDialogsOpenByOperatorID(operatorID)

			// lock dialogs
			dialogLock.Lock()
			dialogs = append(dialogs, d...)
			dialogLock.Unlock()
		}(o.ID)
	}

	wg.Wait()
	fmt.Println(dialogs)
}
