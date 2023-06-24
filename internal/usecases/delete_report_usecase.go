package usecases

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"c2d-reports/internal/services"
	"fmt"
	"log"
	"sync"
)

func (u *ReportTmrUsecase) DeleteReport() {
	var wg sync.WaitGroup

	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	repository := repositories.NewReportRepository(db)
	var reports []repositories.Report
	reports, err = repository.FindAll()
	if err != nil {
		log.Fatalf("could not get reports: %v", err)
	}

	for _, r := range reports {
		wg.Add(1)
		go func(dialogID int) {
			defer wg.Done()
			dialog := chat2deskService.GetDialogByID(dialogID)

			if dialog.State == "closed" || dialog.End != "" {
				err := repository.DeleteReportByDialogID(dialog.ID)
				if err != nil {
					log.Fatalf("failed to delete report: %v", err)
				}
				fmt.Printf("Report with dialog_id %d was deleted\n", dialog.ID)
			}
			fmt.Println("nothing to delete")

		}(r.DialogID)
	}

	// wait for all goroutines to finish
	wg.Wait()
}
