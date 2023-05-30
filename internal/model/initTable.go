package model

import (
	"Time_k8s_operator/internal/dao/db"
)

func InitTable() error {
	if err := db.DB.AutoMigrate(&User{}, &TemplateKind{}, &SpaceTemplate{}, &SpaceSpec{}, &Space{}); err != nil {
		return err
	}
	return nil
}
