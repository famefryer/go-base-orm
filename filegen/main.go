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
				matched := regex.FindString(line)
				fmt.Printf("Annotation = %s\n", matched)

				return matched, nil
			}
		}
	}

	return "", nil
}

func main() {
	pattern := `@GormRepository\(\{"[A-Za-z]*",\s?"[A-Za-z]*"\}\)`
	err := filepath.Walk("./model", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		if !info.IsDir() {
			_, err := scanAnnotation("@GormRepository", pattern, path)
			if err != nil {
				return err
			}
		}
		fmt.Printf("=====================================\n")
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
