package annoscanner

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const StructDeclarationPattern = `type [a-zA-Z]* struct {`
const CommentDeclarationPattern = `//.*`
const StructLastLine = "}"

func ScanAnnotation(annotation, annotationPattern, filepath string) (Model, error) {
	structPattern, err := regexp.Compile(StructDeclarationPattern)
	commentPattern, err := regexp.Compile(CommentDeclarationPattern)

	model := Model{
		Name:       "",
		Annotation: "",
		Attributes: make(map[string]ModelAttribute),
	}

	regex, err := regexp.Compile(annotationPattern)
	if err != nil {
		return model, err
	}

	file, err := os.Open(filepath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		//Loop through each line
		startOfModel := false
		lastLineOfModel := StructLastLine
		for scanner.Scan() {
			line := scanner.Text()

			// Start extracting model's attribute
			if startOfModel {
				// Skip the struct declaration line and all comments
				if commentPattern.MatchString(line) {
					continue
				}

				// Extract model name from the first line
				if structPattern.MatchString(line) {
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
				model.Annotation = regex.FindString(line)
				startOfModel = true
			}
		}
	}

	return model, nil
}
