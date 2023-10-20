package main

type GormRepositoryAnnotation struct {
	TableName  string
	PrimaryKey string
	Columns    map[string]ModelAttribute
}
