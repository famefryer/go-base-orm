package repository

import (
	"base-orm/model/submode"
	"fmt"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db         *gorm.DB
	tableName  string
	primaryKey string
}

// FindByPK query object from database with primary key.
// return object and nil if it exists in database
func (r *ProductRepository) FindByPK(id string) (submode.Product, error) {
	var result submode.Product
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)
	return result, tx.Error
}

// FindAll get all records in database
// return list if all records in database and error
func (r *ProductRepository) FindAll() ([]submode.Product, error) {
	var result []submode.Product
	tx := r.db.Table(r.tableName).Find(&result)
	return result, tx.Error
}

// Create insert new record into database
// return nil if create success
func (r *ProductRepository) Create(object submode.Product) error {
	tx := r.db.Table(r.tableName).Create(object)
	return tx.Error
}

// DeleteByPK record in database by primary key
// return nil if delete success
func (r *ProductRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&submode.Product{})
	return tx.Error
}

// UpdateByPK update existing record in database
// return nil if update success
func (r *ProductRepository) UpdateByPK(object submode.Product) error {
	updatesMap := map[string]interface{}{
		"ID":    object.ID,
		"Name":  object.Name,
		"Price": object.Price,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)
	return tx.Error
}

// NewProductRepository create new gorm repository instance for Product
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db, tableName: "product", primaryKey: " id"}
}
