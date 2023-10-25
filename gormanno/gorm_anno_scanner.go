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

const TypeGormDB = "*gorm.DB"

//var DefaultTypes = []string{"string", "bool", "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "float32", "float64", "complex64", "complex128"}

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
		ModelName:    model.Name,
		ModelPackage: model.Package,
		TableName:    gormValue[0],
		PrimaryKey:   gormValue[1],
		Attributes:   model.Attributes,
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
	fmt.Println("Start generating gorm samplerepository....")

	packageName := "repository"

	gormRepos, err := g.scan(modelDir)
	if err != nil {
		return err
	}

	// Gen output directory
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

		// Generate jennifer file
		modelName := fmt.Sprintf("%sRepository", gormRepo.ModelName)
		f := jen.NewFile(packageName)
		f.ImportName("gorm.io/gorm", "")
		f.Type().Id(modelName).Struct(
			jen.Id("db").Id(TypeGormDB),
			jen.Id("tableName").String(),
			jen.Id("primaryKey").String(),
		)
		f.Func().Params(
			jen.Id("r").Id(fmt.Sprintf("*%s", modelName)),
		).Id("GetByPK").Params(
			jen.Id("id").String(),
		).Params(
			jen.Id("model.User"),
			jen.Id("error"),
		).Block(
			jen.Qual("gorm.io/gorm", "gorm").Call(),
		)
		fmt.Printf("%#v", f)
		buff := &bytes.Buffer{}
		err = f.Render(buff)
		if err != nil {
			fmt.Println(err)
		}

		finalCode := buff.String()
		_, err = file.WriteString(finalCode)
	}

	fmt.Println("Finish generating gorm samplerepository....")

	return nil
}
