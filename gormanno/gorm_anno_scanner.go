package gormanno

import (
	"base-orm/annoscanner"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const GormRepoAnnoPattern = `@GormRepository\("[A-Za-z]*",\s?"[A-Za-z]*"\)`
const GormRepoParamPattern = `"[A-Za-z]*",\s?"[A-Za-z]*"`
const GormRepoAnno = "@GormRepository"

type GormAnnotationScanner struct {
}

func (g *GormAnnotationScanner) genGormRepository(annotation string, attributes map[string]annoscanner.ModelAttribute) error {
	regex, err := regexp.Compile(GormRepoParamPattern)
	matched := regex.FindString(annotation)
	if err != nil {
		return err
	}

	gormValue := strings.Split(matched, ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		TableName:  gormValue[0],
		PrimaryKey: gormValue[1],
		Columns:    attributes,
	}

	fmt.Println(gormAnnoRepo)

	return nil
}

func (g *GormAnnotationScanner) Scan(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			anno, attributes, err := annoscanner.ScanAnnotation(GormRepoAnno, GormRepoAnnoPattern, path)
			if err != nil {
				return err
			}
			if anno == "" {
				// skipped file with no annotation
				return nil
			}

			fmt.Printf("filename: %s\n", info.Name())
			err = g.genGormRepository(anno, attributes)
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

	return nil
}
