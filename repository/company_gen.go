package repository

import (
	"base-orm/model"
	"fmt"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db         *gorm.DB
	tableName  string
	primaryKey string
}

func (r *CompanyRepository) GetByPK(id string) (model.Company, error) {
	var result model.Company
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)
	return result, tx.Error
}

func (r *CompanyRepository) Create(object model.Company) error {
	tx := r.db.Table(r.tableName).Create(object)
	return tx.Error
}

func (r *CompanyRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})
	return tx.Error
}

func (r *CompanyRepository) UpdateByPK(object model.Company) error {
	updatesMap := map[string]interface{}{
		"Branch": object.Branch,
		"Name":   object.Name,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)
	return tx.Error
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return CompanyRepository{db: db, tableName: "company", primaryKey: " name"}
}
