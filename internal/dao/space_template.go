package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
)

func FindOneTemplateByName(name string) (space_template model.SpaceTemplate, ok bool) {
	db.DB.Where("name = ?", name).First(&space_template)
	return space_template, space_template == model.SpaceTemplate{}
}

func InsertSpaceTemplate(space_template *model.SpaceTemplate) (st_id uint32, err error) {
	query := db.DB.Create(space_template)
	err = query.Error
	if err != nil {
		return 0, err
	}
	var s model.SpaceTemplate
	query.First(&s)
	return s.Id, nil
}
