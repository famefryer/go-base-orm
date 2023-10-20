package gormanno

import (
	"base-orm/annoscanner"
	"fmt"
	"regexp"
	"strings"
)

const GormRepoAnnoPattern = `@GormRepository\("[A-Za-z]*",\s?"[A-Za-z]*"\)`
const GormRepoParamPattern = `"[A-Za-z]*",\s?"[A-Za-z]*"`

func GenGormRepository(annotation string) error {
	regex, err := regexp.Compile(GormRepoParamPattern)
	matched := regex.FindString(annotation)
	if err != nil {
		return err
	}

	columnMap := make(map[string]annoscanner.ModelAttribute)
	gormValue := strings.Split(matched, ",")
	gormAnnoRepo := GormRepositoryAnnotation{
		TableName:  gormValue[0],
		PrimaryKey: gormValue[1],
		Columns:    columnMap,
	}

	fmt.Println(gormAnnoRepo)

	return nil
}
