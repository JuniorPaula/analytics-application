package usecases

import (
	"c2d-reports/internal/services"
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

	for _, o := range operators {

		chat2deskService.GetAllDialogsOpenByOperatorID(o.ID)
	}
}
