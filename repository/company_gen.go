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

// GetByPK query object from database with primary key.
// return object and nil if it exists in database
func (r *CompanyRepository) GetByPK(id string) (model.Company, error) {
	var result model.Company
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)
	return result, tx.Error
}

// Create create new record in database
// return nil if create success
func (r *CompanyRepository) Create(object model.Company) error {
	tx := r.db.Table(r.tableName).Create(object)
	return tx.Error
}

// DeleteByPK record in database by primary key
// return nil if delete success
func (r *CompanyRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})
	return tx.Error
}

// UpdateByPK update existing record in database
// return nil if update success
func (r *CompanyRepository) UpdateByPK(object model.Company) error {
	updatesMap := map[string]interface{}{
		"Branch": object.Branch,
		"Name":   object.Name,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)
	return tx.Error
}

// NewCompanyRepository create new gorm repository instance for Company
func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db, tableName: "company", primaryKey: " name"}
}
