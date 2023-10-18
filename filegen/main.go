package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func scanAnnotation(annotation, annotationPattern, filepath string) (string, error) {
	regex, err := regexp.Compile(annotationPattern)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filepath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		//Loop through each line
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "//") && strings.Contains(line, annotation) {
				result := regex.FindString(line)
				//fmt.Printf("Annotation = %s\n", result)

				return result, nil
			}
		}
	}

	return "", nil
}

func genGormRepository(annotation string) error {
	pattern := `"[A-Za-z]*",\s?"[A-Za-z]*"`
	regex, err := regexp.Compile(pattern)
	matched := regex.FindString(annotation)
	if err != nil {
		return err
	}

	gormValue := strings.Split(matched, ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		tableName:  gormValue[0],
		primaryKey: gormValue[1],
	}

	fmt.Println(gormAnnoRepo)

	return nil
}

func main() {
	pattern := `@GormRepository\("[A-Za-z]*",\s?"[A-Za-z]*"\)`
	err := filepath.Walk("./model", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			anno, err := scanAnnotation("@GormRepository", pattern, path)
			if err != nil {
				return err
			}
			if anno == "" {
				// skipped file with no annotation
				return nil
			}

			fmt.Printf("filename: %s\n", info.Name())
			err = genGormRepository(anno)
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
