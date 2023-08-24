package entity

import "errors"

type Company struct {
	ID           int64  `json:"id,omitempty"`
	CompanyID    int    `json:"company_id,omitempty"`
	CompanyToken string `json:"company_token,omitempty"`
	CompanyName  string `json:"company_name,omitempty"`
	EmailAdmin   string `json:"email_admin,omitempty"`
}

// Validate validates the company
func (c *Company) Validate() error {
	if c.CompanyID == 0 {
		return errors.New("company_id is empty")
	}
	if c.CompanyToken == "" {
		return errors.New("company_token is empty")
	}
	if c.CompanyName == "" {
		return errors.New("company_name is empty")
	}
	if c.EmailAdmin == "" {
		return errors.New("email_admin is empty")
	}
	return nil
}
