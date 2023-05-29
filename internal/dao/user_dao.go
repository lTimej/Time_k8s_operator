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

func FindOneUserByUsername(username string) bool {
	var user model.User
	db.DB.Where("username = ?", username).First(&user)
	return model.User{} == user
}

func FindOneUserByUsernameAndPassword(username, password string)(*model.User,bool) {
	var user model.User
	db.DB.Where("username = ? AND password = ?", username, password).First(&user)
	return &user,model.User{} == user
}
