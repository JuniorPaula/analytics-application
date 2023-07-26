package companies

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/entity"
	"c2d-reports/internal/repositories"
)

func CreateCompanyUsecase(company entity.Company) (entity.Company, error) {
	if err := company.Validate(); err != nil {
		return entity.Company{}, err
	}

	db, err := database.Connect()
	if err != nil {
		return entity.Company{}, err
	}

	repository := repositories.NewCompanyRepository(db)
	company, err = repository.CreateCompany(company)
	if err != nil {
		return entity.Company{}, err
	}

	return entity.Company{}, nil
}
