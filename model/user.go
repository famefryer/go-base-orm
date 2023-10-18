package model

// @GormRepository({"tableName", "primarykey"})
type User struct {
	Username string `gorm:"primarykey"`
	Name     string
	Age      int
}
