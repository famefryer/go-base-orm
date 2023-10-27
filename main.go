package main

import (
	repository "base-orm/gen"
	"base-orm/gormanno"
	"fmt"
)

func main() {
	// Scan gorm annotation
	gormAnnoGen := gormanno.GormAnnotationScanner{}
	err := gormAnnoGen.Execute("./model", "./gen")
	if err != nil {
		fmt.Println(err)
	}
	repository.CompanyRepository{}
}
