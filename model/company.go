package model

// @GormRepository("company", "name")
type Company struct {
	Name   string
	Branch CompanyBranch
}
