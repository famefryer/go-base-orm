package main

import (
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
}
