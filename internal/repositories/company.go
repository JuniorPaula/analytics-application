package repositories

import (
	"c2d-reports/internal/entity"
	"database/sql"
)

type CompanyRepository struct {
	Db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{
		Db: db,
	}
}

func (c *CompanyRepository) CreateCompany(company entity.Company) (entity.Company, error) {
	statment, err := c.Db.Prepare("INSERT INTO companies (company_id, company_token, company_name, email_admin) VALUES (?, ?, ?, ?) RETURNING id")
	if err != nil {
		return entity.Company{}, err
	}

	var id int
	err = statment.QueryRow(company.CompanyID, company.CompanyToken, company.CompanyName, company.EmailAdmin).Scan(&id)
	if err != nil {
		return entity.Company{}, err
	}

	company.ID = id

	return company, nil
}
