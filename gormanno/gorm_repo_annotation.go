package gormanno

import (
	"base-orm/annoscanner"
)

type GormRepositoryAnnotation struct {
	ModelName  string
	TableName  string
	PrimaryKey string
	Attributes map[string]annoscanner.ModelAttribute
}
