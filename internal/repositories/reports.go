package repositories

import (
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

type Report struct {
	ID            int64  `json:"id,omitempty"`
	OperatorName  string `json:"operator_name"`
	OperatorID    int    `json:"operator_id"`
	DialogID      int    `json:"dialog_id"`
	TMRInSeconds  int    `json:"tmr_in_seconds"`
	OpenedDialogs int    `json:"opened_dialogs"`
	Client        string `json:"client"`
	StatusTAG     string `json:"status_tag"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db}
}

func (r *ReportRepository) CreateOrUpdate(report Report) (uint64, error) {
	var id uint64
	err := r.db.QueryRow(`
		INSERT INTO reports (operator_name, operator_id, dialog_id, tmr_in_seconds, opened_dialogs, client, status_tag)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			operator_name = VALUES(operator_name),
			operator_id = VALUES(operator_id),
			dialog_id = VALUES(dialog_id),
			tmr_in_seconds = VALUES(tmr_in_seconds),
			opened_dialogs = VALUES(opened_dialogs),
			client = VALUES(client),
			status_tag = VALUES(status_tag)
	`, report.OperatorName, report.OperatorID, report.DialogID, report.TMRInSeconds, report.OpenedDialogs, report.Client, report.StatusTAG).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReportRepository) FindAll() ([]Report, error) {
	rows, err := r.db.Query(`
		SELECT id, operator_name, operator_id, dialog_id, tmr_in_seconds, opened_dialogs, client, status_tag, created_at, updated_at
		FROM reports
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var report Report
		err := rows.Scan(&report.ID, &report.OperatorName, &report.OperatorID, &report.DialogID, &report.TMRInSeconds, &report.OpenedDialogs, &report.Client, &report.StatusTAG, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}
