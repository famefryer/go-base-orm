package model

type User struct {
	Username string `gorm:"primarykey"`
	Name     string
	Age      int
}
