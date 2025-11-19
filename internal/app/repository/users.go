package repository

import (
	"LAB1/internal/app/ds"
)

func (r *Repository) CreateUser(u *ds.Users) error {
	return r.db.Create(u).Error
}

func (r *Repository) GetUserByLogin(login string) (ds.Users, error) {
	var u ds.Users
	err := r.db.Where("login = ?", login).First(&u).Error
	return u, err
}

func (r *Repository) GetUserByID(id uint) (ds.Users, error) {
	var u ds.Users
	err := r.db.First(&u, id).Error
	return u, err
}
