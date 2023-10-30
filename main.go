package main

import (
	"base-orm/gormanno"
	"fmt"
)

func main() {
	// Scan gorm annotation
	gormAnnoGen := gormanno.GormAnnotationScanner{}
	err := gormAnnoGen.Execute("base-orm", "./model", "./repository")
	if err != nil {
		fmt.Println(err)
	}
}
