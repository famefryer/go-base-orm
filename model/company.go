package model

// @GormRepository("company", "name")
type Company struct {
	Name   string `gorm:"primarykey"`
	Branch CompanyBranch
}
