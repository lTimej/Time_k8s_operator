package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
)

func FindAllTemplateKind() (kinds []model.TemplateKind) {
	db.DB.Find(&kinds)
	return
}
