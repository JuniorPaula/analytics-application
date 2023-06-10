package usecases

import (
	"c2d-reports/internal/services"
	"fmt"
)

type ReportTmrUsecase struct {
	CompanyToken string
}

func (u *ReportTmrUsecase) LoadTMR() {

	// calls chat2desk api
	chat2deskService := services.Chat2DeskService{
		Token: u.CompanyToken,
	}
	operators := chat2deskService.GetOperators()

	fmt.Println(operators)
}
