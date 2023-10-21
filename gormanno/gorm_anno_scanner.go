package gormanno

import (
	"base-orm/annoscanner"
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
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

func (g *GormAnnotationScanner) genGormRepository(model annoscanner.Model) (GormRepositoryAnnotation, error) {
	regex, err := regexp.Compile(GormRepoParamPattern)
	matched := regex.FindString(model.Annotation)
	if err != nil {
		return GormRepositoryAnnotation{}, err
	}

	gormValue := strings.Split(strings.ReplaceAll(matched, "\"", ""), ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		ModelName:  model.Name,
		TableName:  gormValue[0],
		PrimaryKey: gormValue[1],
		Columns:    model.Attributes,
	}

	return gormAnnoRepo, nil
}

func (g *GormAnnotationScanner) scan(dir string) ([]GormRepositoryAnnotation, error) {
	var gormRepos []GormRepositoryAnnotation
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			model, err := annoscanner.ScanAnnotation(GormRepoAnno, GormRepoAnnoPattern, path)
			if err != nil {
				return err
			}
			if model.Annotation == "" {
				// skipped file with no annotation
				return nil
			}

			gormRepo, err := g.genGormRepository(model)
			if err != nil {
				return err
			}

			gormRepos = append(gormRepos, gormRepo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gormRepos, nil
}

func (g *GormAnnotationScanner) Execute(modelDir, outputDir string) error {
	fmt.Println("Start generating gorm repository....")

	packageName := "repository"

	gormRepos, err := g.scan(modelDir)
	if err != nil {
		return err
	}

	// Gen output path
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0777)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Gen output files
	for _, gormRepo := range gormRepos {
		filename := fmt.Sprintf("%s_gen", strings.ToLower(gormRepo.ModelName))
		file, err := os.Create(fmt.Sprintf("./gen/%s.go", filename))

		if err != nil {
			fmt.Println(err)
		}

		f := jen.NewFile(packageName)

		f.Func().Id("main").Params().Block(
			jen.Qual("a.b/c", "Foo").Call(),
		)

		buff := &bytes.Buffer{}
		err = f.Render(buff)
		if err != nil {
			fmt.Println(err)
		}
		_, err = file.WriteString(buff.String())
	}

	fmt.Println("Finish generating gorm repository....")

	return nil
}
