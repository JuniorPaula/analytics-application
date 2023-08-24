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
	statment, err := c.Db.Prepare("INSERT INTO companies (company_id, company_token, company_name, email_admin) VALUES (?, ?, ?, ?)")
	if err != nil {
		return entity.Company{}, err
	}

	result, err := statment.Exec(company.CompanyID, company.CompanyToken, company.CompanyName, company.EmailAdmin)
	if err != nil {
		return entity.Company{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.Company{}, err
	}

	company.ID = id

	return company, nil
}

func (c *CompanyRepository) GetAllCompanies() ([]entity.Company, error) {
	rows, err := c.Db.Query("SELECT id, company_id, company_token, company_name, email_admin FROM companies")
	if err != nil {
		return nil, err
	}

	var companies []entity.Company

	for rows.Next() {
		var company entity.Company

		err = rows.Scan(&company.ID, &company.CompanyID, &company.CompanyToken, &company.CompanyName, &company.EmailAdmin)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}

func (c *CompanyRepository) GetCompanyByID(ID int64) (entity.Company, error) {
	var company entity.Company

	row := c.Db.QueryRow("SELECT id, company_id, company_token, company_name, email_admin FROM companies WHERE id = ?", ID)

	err := row.Scan(&company.ID, &company.CompanyID, &company.CompanyToken, &company.CompanyName, &company.EmailAdmin)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func (c *CompanyRepository) UpdateCompany(ID int64, company entity.Company) (entity.Company, error) {
	statment, err := c.Db.Prepare("UPDATE companies SET company_id = ?, company_token = ?, company_name = ?, email_admin = ? WHERE id = ?")
	if err != nil {
		return entity.Company{}, err
	}

	_, err = statment.Exec(company.CompanyID, company.CompanyToken, company.CompanyName, company.EmailAdmin, ID)
	if err != nil {
		return entity.Company{}, err
	}

	company.ID = ID

	return company, nil
}
