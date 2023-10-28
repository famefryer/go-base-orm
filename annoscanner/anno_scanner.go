package annoscanner

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const StructDeclarationPattern = `type [a-zA-Z]* struct {`
const CommentDeclarationPattern = `//.*`
const PackagePattern = `package .*`
const StructLastLine = "}"

func ScanAnnotation(annotation, annotationPattern, filepath string) (Model, error) {
	structRegex, err := regexp.Compile(StructDeclarationPattern)
	commentRegex, err := regexp.Compile(CommentDeclarationPattern)
	packageRegex, err := regexp.Compile(PackagePattern)

	model := Model{
		Name:       "",
		Annotation: "",
		Attributes: make(map[string]ModelAttribute),
	}

	annoRegex, err := regexp.Compile(annotationPattern)
	if err != nil {
		return model, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return model, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	//Loop through each line
	startOfModel := false
	lastLineOfModel := StructLastLine
	for scanner.Scan() {
		line := scanner.Text()

		// Extract package name
		if packageRegex.MatchString(line) {
			packages := strings.Fields(line)
			model.Package = packages[1]
		}

		// Start extracting model's attribute
		if startOfModel {
			// Skip the struct declaration line and all comments
			if commentRegex.MatchString(line) {
				continue
			}

			// Extract model name from the first line
			if structRegex.MatchString(line) {
				structDC := strings.Fields(strings.TrimSpace(line))
				model.Name = structDC[1]
				continue
			}

			// Stop finding attribute when reaching the end of struct
			if strings.TrimSpace(line) == lastLineOfModel {
				return model, nil
			}

			// Extract model's attribute to this sample format {name:username, dataType:string}
			attribute := strings.Fields(strings.TrimSpace(line))
			model.Attributes[attribute[0]] = MakeModelAttribute(attribute[0], attribute[1])
		}

		// Toggle startOfModel when it detects annotation
		if !startOfModel && strings.Contains(line, "//") && strings.Contains(line, annotation) {
			model.Annotation = annoRegex.FindString(line)
			startOfModel = true
		}
	}

	return model, nil
}
