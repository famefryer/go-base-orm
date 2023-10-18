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

func (r *UserRepository) GetByPK(id string) (model.User, error) {
	var result model.User
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).First(&result)

	return result, tx.Error
}

func (r *UserRepository) Create(object model.User) error {
	tx := r.db.Table(r.tableName).Create(object)

	return tx.Error
}

func (r *UserRepository) DeleteByPK(id string) error {
	query := fmt.Sprintf("%s = ?", r.primaryKey)
	tx := r.db.Table(r.tableName).Where(query, id).Delete(&model.User{})

	return tx.Error
}

func (r *UserRepository) UpdateByPK(object model.User) error {
	updatesMap := map[string]interface{}{
		"username": object.Username,
		"name":     object.Name,
		"age":      object.Age,
	}
	tx := r.db.Table(r.tableName).Updates(updatesMap)

	return tx.Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db:         db,
		tableName:  "user",
		primaryKey: "username",
	}
}
