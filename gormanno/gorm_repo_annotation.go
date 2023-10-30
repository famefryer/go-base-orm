package gormanno

import (
	"base-orm/annoscanner"
)

type GormRepositoryAnnotation struct {
	ModelImportPackagePath string
	ModelPackage           string
	ModelName              string
	TableName              string
	PrimaryKey             string
	Attributes             map[string]annoscanner.ModelAttribute
}
