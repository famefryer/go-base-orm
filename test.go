package main

import (
	"base-orm/repository"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dns := "root:@tcp(127.0.0.1:3306)/playground?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		fmt.Printf("error : %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	user, err := userRepo.GetByID("username1")
	fmt.Printf("User : %v", user)
}
