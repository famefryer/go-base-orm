package model

// @GormRepository("company", "name")

type Company struct {
	name  string
	users []User
}
