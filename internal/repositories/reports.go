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
	OperatorID    int64  `json:"operator_id"`
	DialogID      int64  `json:"dialog_id"`
	TMRInSeconds  int64  `json:"tmr_in_seconds"`
	OpenedDialogs int64  `json:"opened_dialogs"`
	Client        string `json:"client"`
	Status        string `json:"status"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db}
}

func (r *ReportRepository) CreateOrUpdate(report Report) (uint64, error) {
	var id uint64
	err := r.db.QueryRow(`
		INSERT INTO reports (operator_name, operator_id, dialog_id, tmr_in_seconds, opened_dialogs, client, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			operator_name = VALUES(operator_name),
			operator_id = VALUES(operator_id),
			dialog_id = VALUES(dialog_id),
			tmr_in_seconds = VALUES(tmr_in_seconds),
			opened_dialogs = VALUES(opened_dialogs),
			client = VALUES(client),
			status = VALUES(status)
	`, report.OperatorName, report.OperatorID, report.DialogID, report.TMRInSeconds, report.OpenedDialogs, report.Client, report.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
