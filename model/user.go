package model

// @GormRepository("user", "username")
type User struct {
	Username string `gorm:"primarykey"`
	// dsf
	Name string
	Age  int
}
