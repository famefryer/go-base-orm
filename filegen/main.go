package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func scanAnnotation(annotation, annotationPattern, filepath string) (string, map[string]ModelAttribute, error) {
	structPattern, err := regexp.Compile(`type [a-zA-Z]* struct {`)
	commentPattern, err := regexp.Compile(`//.*`)
	modelAttrs := make(map[string]ModelAttribute)

	regex, err := regexp.Compile(annotationPattern)
	if err != nil {
		return "", modelAttrs, err
	}

	file, err := os.Open(filepath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		//Loop through each line
		startOfModel := false
		lastLineOfModel := "}"
		annotationString := ""
		for scanner.Scan() {
			line := scanner.Text()

			if startOfModel {
				if structPattern.MatchString(line) || commentPattern.MatchString(line) {
					continue
				}

				if strings.TrimSpace(line) == lastLineOfModel {
					fmt.Println("Fame")
					fmt.Printf("%v\n", modelAttrs)
					fmt.Println("-------")
					return annotationString, modelAttrs, nil
				}

				attribute := strings.Split(strings.TrimSpace(scanner.Text()), " ")
				modelAttrs[attribute[0]] = MakeModelAttribute(attribute[0], attribute[1])
			}

			if !startOfModel && strings.Contains(line, "//") && strings.Contains(line, annotation) {
				annotationString = regex.FindString(line)
				startOfModel = true
			}
		}
	}

	return "", modelAttrs, nil
}

func genGormRepository(annotation string) error {
	pattern := `"[A-Za-z]*",\s?"[A-Za-z]*"`
	regex, err := regexp.Compile(pattern)
	matched := regex.FindString(annotation)
	if err != nil {
		return err
	}

	columnMap := make(map[string]ModelAttribute)
	gormValue := strings.Split(matched, ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		TableName:  gormValue[0],
		PrimaryKey: gormValue[1],
		Columns:    columnMap,
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
			anno, _, err := scanAnnotation("@GormRepository", pattern, path)
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
