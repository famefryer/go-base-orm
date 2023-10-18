package model

// @GormRepository("user", "username")
type User struct {
	Username string `gorm:"primarykey"`
	Name     string
	Age      int
}
