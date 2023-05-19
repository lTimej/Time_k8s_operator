package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
)

func FindOneUserByEmail(email string) bool {
	var user model.User
	db.DB.Where("email = ?", email).First(&user)
	return model.User{} != user
}
