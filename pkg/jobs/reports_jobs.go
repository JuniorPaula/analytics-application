package jobs

import (
	"c2d-reports/internal/usecases/companies"
	usecases "c2d-reports/internal/usecases/reports"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Schedule struct {
	Cron *cron.Cron
}

// NewSchedule returns a new Schedule
func NewSchedule() *Schedule {
	return &Schedule{
		Cron: cron.New(),
	}
}

// ScheduleCalculateReport schedules the job to calculate the report
func (s *Schedule) ScheduleCalculateReport() {
	cron := s.Cron

	companies, err := companies.GetAllCompaniesUsecase()
	if err != nil {
		panic(err)
	}

	for _, c := range companies {
		// create a copy of the company
		company := c

		// create the usecase
		uc := usecases.ReportTmrUsecase{
			CompanyToken: company.CompanyToken,
		}
		// add the job to the cron
		_, err = cron.AddFunc("*/1 8-18 * * *", func() {
			fmt.Printf("start schedule to computation report chats from company: [%s]\n", company.CompanyName)
			uc.LoadTMR()
		})
		if err != nil {
			panic(err)
		}
	}
	cron.Start()

}

// ScheduleDeleteReport schedules the job to delete the report
func (s *Schedule) ScheduleDeleteReport() {
	cron := s.Cron

	companies, err := companies.GetAllCompaniesUsecase()
	if err != nil {
		panic(err)
	}

	for _, c := range companies {
		// create a copy of the company
		company := c

		// create the usecase
		uc := usecases.ReportTmrUsecase{
			CompanyToken: company.CompanyToken,
		}

		// add the job to the cron
		_, err = cron.AddFunc("*/2 6-20 * * *", func() {
			fmt.Printf("start schedule to delete report from company: [%s]\n", company.CompanyName)
			uc.DeleteReport()
		})
		if err != nil {
			panic(err)
		}
		cron.Start()
	}

}
