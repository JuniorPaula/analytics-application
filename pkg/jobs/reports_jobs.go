package jobs

import (
	"c2d-reports/internal/config"
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
	c := s.Cron
	uc := usecases.ReportTmrUsecase{
		CompanyToken: config.CompanyToken,
	}
	_, err := c.AddFunc("*/1 8-18 * * *", func() {
		fmt.Println("start schedule to computation report chats")
		uc.LoadTMR()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

// ScheduleDeleteReport schedules the job to delete the report
func (s *Schedule) ScheduleDeleteReport() {
	cron := s.Cron

	companies, err := companies.GetAllCompaniesUsecase()
	if err != nil {
		panic(err)
	}

	for _, c := range companies {
		uc := usecases.ReportTmrUsecase{
			CompanyToken: c.CompanyToken,
		}
		_, err = cron.AddFunc("*/2 6-20 * * *", func() {
			fmt.Println("start schedule to delete report")
			uc.DeleteReport()
		})
		if err != nil {
			panic(err)
		}
		cron.Start()
	}

}
