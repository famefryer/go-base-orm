package main

import (
	"base-orm/gormanno"
	"fmt"
)

func main() {
	gormAnnoGen := gormanno.GormAnnotationScanner{}
	err := gormAnnoGen.Scan("./model")
	if err != nil {
		fmt.Println(err)
	}
}
