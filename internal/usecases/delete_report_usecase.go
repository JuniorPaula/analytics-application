package usecases

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"c2d-reports/internal/services"
	"fmt"
	"log"
)

func (u *ReportTmrUsecase) DeleteReport() {

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
		dialog := chat2deskService.GetDialogByID(r.DialogID)
		fmt.Println(dialog)
	}

}
