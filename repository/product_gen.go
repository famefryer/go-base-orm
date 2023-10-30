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

func (r *ProductRepository) GetByPK(id string) (submode.Product, error) {
	var result submode.Product
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)
	return result, tx.Error
}

func (r *ProductRepository) Create(object submode.Product) error {
	tx := r.db.Table(r.tableName).Create(object)
	return tx.Error
}

func (r *ProductRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})
	return tx.Error
}

func (r *ProductRepository) UpdateByPK(object submode.Product) error {
	updatesMap := map[string]interface{}{
		"ID":    object.ID,
		"Name":  object.Name,
		"Price": object.Price,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)
	return tx.Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db: db, tableName: "product", primaryKey: " id"}
}
