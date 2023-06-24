package usecases

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/repositories"
	"fmt"
	"log"
)

func (u *ReportTmrUsecase) DeleteReport() {

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

	fmt.Println(reports)
}
