package main

import (
	"base-orm/annoscanner"
	"base-orm/gormanno"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := filepath.Walk("./model", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			anno, _, err := annoscanner.ScanAnnotation("@GormRepository", gormanno.GormRepoAnnoPattern, path)
			if err != nil {
				return err
			}
			if anno == "" {
				// skipped file with no annotation
				return nil
			}

			fmt.Printf("filename: %s\n", info.Name())
			err = gormanno.GenGormRepository(anno)
			if err != nil {
				return err
			}
			fmt.Println("=============================")
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
