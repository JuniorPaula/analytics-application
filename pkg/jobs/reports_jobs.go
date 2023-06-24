package jobs

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/usecases"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Schedule struct {
	Cron *cron.Cron
}

func NewSchedule() *Schedule {
	return &Schedule{
		Cron: cron.New(),
	}
}

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

func (s *Schedule) ScheduleDeleteReport() {
	c := s.Cron
	uc := usecases.ReportTmrUsecase{
		CompanyToken: config.CompanyToken,
	}
	_, err := c.AddFunc("*/2 6-20 * * *", func() {
		fmt.Println("start schedule to delete report")
		uc.DeleteReport()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
