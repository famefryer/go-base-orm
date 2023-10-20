package gormanno

import (
	"base-orm/annoscanner"
)

type GormRepositoryAnnotation struct {
	TableName  string
	PrimaryKey string
	Columns    map[string]annoscanner.ModelAttribute
}
