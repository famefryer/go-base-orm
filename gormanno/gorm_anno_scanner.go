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
	projectName     string
	modelDirectory  string
	outputDirectory string
}

func (g *GormAnnotationScanner) genGormRepository(model annoscanner.Model) (GormRepositoryAnnotation, error) {
	regex, err := regexp.Compile(GormRepoParamPattern)
	matched := regex.FindString(model.Annotation)
	if err != nil {
		return GormRepositoryAnnotation{}, err
	}

	gormValue := strings.Split(strings.ReplaceAll(matched, "\"", ""), ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		ModelImportPackagePath: model.ImportPackagePath,
		ModelName:              model.Name,
		ModelPackage:           model.Package,
		TableName:              gormValue[0],
		PrimaryKey:             gormValue[1],
		Attributes:             model.Attributes,
	}

	return gormAnnoRepo, nil
}

func (g *GormAnnotationScanner) scan() ([]GormRepositoryAnnotation, error) {
	var gormRepos []GormRepositoryAnnotation
	err := filepath.Walk(g.modelDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			model, err := annoscanner.ScanAnnotation(g.projectName, GormRepoAnno, GormRepoAnnoPattern, path)
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

func (g *GormAnnotationScanner) genOutputFile(gormRepos []GormRepositoryAnnotation) error {
	packageName := "repository"

	for _, gormRepo := range gormRepos {
		filename := fmt.Sprintf("%s_gen", strings.ToLower(gormRepo.ModelName))
		file, err := os.Create(fmt.Sprintf("./%s/%s.go", g.outputDirectory, filename))

		if err != nil {
			fmt.Println(err)
		}

		// Generate jennifer file
		repoName := fmt.Sprintf("%sRepository", gormRepo.ModelName)
		modelType := fmt.Sprintf("%s.%s", gormRepo.ModelPackage, gormRepo.ModelName)
		queryByPK := "fmt.Sprintf(\"%s = ?\", r.primaryKey)"
		attrMap := jen.Dict{}
		for _, attr := range gormRepo.Attributes {
			attrMap[jen.Lit(attr.Name)] = jen.Id(fmt.Sprintf("object.%s", attr.Name))
		}

		f := jen.NewFile(packageName)

		f.Id(fmt.Sprintf("import (\n\t\"gorm.io/gorm\"\n\t\"fmt\"\n\t\"%s\"\n)", gormRepo.ModelImportPackagePath))

		f.Type().Id(repoName).Struct(
			jen.Id("db").Id(TypeGormDB),
			jen.Id("tableName").String(),
			jen.Id("primaryKey").String(),
		)

		// function GetByPK
		f.Func().Params(
			jen.Id("r").Id(fmt.Sprintf("*%s", repoName)),
		).Id("GetByPK").Params(
			jen.Id("id").String(),
		).Params(
			jen.Id(modelType),
			jen.Id("error"),
		).Block(
			jen.Var().Id("result").Id(modelType),
			jen.Id("query").Op(":=").Id(queryByPK),
			jen.Id("tx").Op(":=").Id("r.db.Table(r.tableName).Where(query, id).First(&result)"),
			jen.Return(jen.Id("result, tx.Error")),
		).Line()

		// function Create
		f.Func().Params(
			jen.Id("r").Id(fmt.Sprintf("*%s", repoName)),
		).Id("Create").Params(
			jen.Id("object").Id(modelType),
		).Params(
			jen.Id("error"),
		).Block(
			jen.Id("tx").Op(":=").Id("r.db.Table(r.tableName).Create(object)"),
			jen.Return(jen.Id("tx.Error")),
		).Line()

		// function DeleteByPK
		f.Func().Params(
			jen.Id("r").Id(fmt.Sprintf("*%s", repoName)),
		).Id("DeleteByPK").Params(
			jen.Id("id").String(),
		).Params(
			jen.Id("error"),
		).Block(
			jen.Id("query").Op(":=").Id(queryByPK),
			jen.Id("tx").Op(":=").Id("r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})"),
			jen.Return(jen.Id("tx.Error")),
		).Line()

		// function UpdateByPK
		f.Func().Params(
			jen.Id("r").Id(fmt.Sprintf("*%s", repoName)),
		).Id("UpdateByPK").Params(
			jen.Id("object").Id(modelType),
		).Params(
			jen.Id("error"),
		).Block(
			jen.Id("updatesMap").Op(":=").Map(jen.String()).Interface().Values(attrMap),
			jen.Id("tx").Op(":=").Id("r.db.Table(r.tableName).Updates(updatesMap)"),
			jen.Return(jen.Id("tx.Error")),
		).Line()

		// function NewGormRepository
		f.Func().Id(fmt.Sprintf("New%s", repoName)).Params(
			jen.Id("db").Id(TypeGormDB),
		).Params(
			jen.Id(repoName),
		).Block(
			jen.Return(jen.Id(fmt.Sprintf("%s{db: db, tableName: \"%s\", primaryKey: \"%s\"}", repoName, gormRepo.TableName, gormRepo.PrimaryKey))),
		).Line()

		buff := &bytes.Buffer{}
		err = f.Render(buff)
		if err != nil {
			fmt.Println(err)
		}

		finalCode := buff.String()
		_, err = file.WriteString(finalCode)
	}

	return nil
}

func (g *GormAnnotationScanner) Execute(projectName, modelDir, outputDir string) error {
	fmt.Println("Start generating gorm repository....")

	g.projectName = projectName
	g.modelDirectory = modelDir
	g.outputDirectory = outputDir

	gormRepos, err := g.scan()
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
	err = g.genOutputFile(gormRepos)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Finish generating gorm repository....")

	return nil
}
