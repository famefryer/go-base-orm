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

func (g *GormAnnotationScanner) genGormRepository(annotation string, attributes map[string]annoscanner.ModelAttribute) (GormRepositoryAnnotation, error) {
	regex, err := regexp.Compile(GormRepoParamPattern)
	matched := regex.FindString(annotation)
	if err != nil {
		return GormRepositoryAnnotation{}, err
	}

	gormValue := strings.Split(strings.ReplaceAll(matched, "\"", ""), ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		TableName:  gormValue[0],
		PrimaryKey: gormValue[1],
		Columns:    attributes,
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
			anno, attributes, err := annoscanner.ScanAnnotation(GormRepoAnno, GormRepoAnnoPattern, path)
			if err != nil {
				return err
			}
			if anno == "" {
				// skipped file with no annotation
				return nil
			}

			gormRepo, err := g.genGormRepository(anno, attributes)
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
		filename := fmt.Sprintf("%s_gen", gormRepo.TableName)
		file, err := os.Create(fmt.Sprintf("./gen/%s.go", filename))

		if err != nil {
			fmt.Println(err)
		}
		f := jen.NewFilePath("./gen/repository")

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
