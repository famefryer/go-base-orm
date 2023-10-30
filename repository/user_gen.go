package repository

import (
	"base-orm/model"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db         *gorm.DB
	tableName  string
	primaryKey string
}

// FindByPK query object from database with primary key.
// return object and nil if it exists in database
func (r *UserRepository) FindByPK(id string) (model.User, error) {
	var result model.User
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)
	return result, tx.Error
}

// FindAll get all records in database
// return list if all records in database and error
func (r *UserRepository) FindAll() ([]model.User, error) {
	var result []model.User
	tx := r.db.Table(r.tableName).Find(&result)
	return result, tx.Error
}

// Create insert new record into database
// return nil if create success
func (r *UserRepository) Create(object model.User) error {
	tx := r.db.Table(r.tableName).Create(object)
	return tx.Error
}

// DeleteByPK record in database by primary key
// return nil if delete success
func (r *UserRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})
	return tx.Error
}

// UpdateByPK update existing record in database
// return nil if update success
func (r *UserRepository) UpdateByPK(object model.User) error {
	updatesMap := map[string]interface{}{
		"Age":         object.Age,
		"CompanyName": object.CompanyName,
		"Name":        object.Name,
		"Username":    object.Username,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)
	return tx.Error
}

// NewUserRepository create new gorm repository instance for User
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db, tableName: "user", primaryKey: " username"}
}
