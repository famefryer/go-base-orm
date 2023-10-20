package annoscanner

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const StructDeclarationPattern = `type [a-zA-Z]* struct {`
const CommentDeclarationPattern = `type [a-zA-Z]* struct {`
const StructLastLine = "}"

func ScanAnnotation(annotation, annotationPattern, filepath string) (string, map[string]ModelAttribute, error) {
	structPattern, err := regexp.Compile(StructDeclarationPattern)
	commentPattern, err := regexp.Compile(CommentDeclarationPattern)
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
		lastLineOfModel := StructLastLine
		annotationString := ""
		for scanner.Scan() {
			line := scanner.Text()

			// Start extracting model's attribute
			if startOfModel {
				// Skip the struct declaration line and all comments
				if structPattern.MatchString(line) || commentPattern.MatchString(line) {
					continue
				}

				// Stop finding attribute when reaching the end of struct
				if strings.TrimSpace(line) == lastLineOfModel {
					return annotationString, modelAttrs, nil
				}

				// Extract model's attribute to this sample format {name:username, dataType:string}
				attribute := strings.Fields(strings.TrimSpace(line))
				modelAttrs[attribute[0]] = MakeModelAttribute(attribute[0], attribute[1])
			}

			// Toggle startOfModel when it detects annotation
			if !startOfModel && strings.Contains(line, "//") && strings.Contains(line, annotation) {
				annotationString = regex.FindString(line)
				startOfModel = true
			}
		}
	}

	return "", modelAttrs, nil
}
