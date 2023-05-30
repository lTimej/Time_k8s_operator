package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
)

const (
	TemplateUsing = iota
	TemplateDeleted
)

func FindAllSpaceTemplateByUsing() (tmps []model.SpaceTemplate) {
	db.DB.Where("status = ?", TemplateUsing).Find(&tmps)
	return
}

func FindAllTemplateKind() (kinds []model.TemplateKind) {
	db.DB.Find(&kinds)
	return
}

func FindAllSpec() (specs []model.SpaceSpec) {
	db.DB.Find(&specs)
	return
}

func FindOneByUserIdAndName(user_id uint32, name string) bool {
	var space model.Space
	db.DB.Where("user_id = ? AND name = ? AND status != ?", user_id, name, model.SpaceStatusDeleted).First(&space)
	return model.Space{} == space
}

func FindCountByUserId(user_id uint32) int64 {
	var count int64
	db.DB.Table("space").Where("user_id").Count(&count)
	return count
}

func InsertSpace(space *model.Space) (space_id uint32, err error) {
	query := db.DB.Create(space)
	err = query.Error
	if err != nil {
		return 0, err
	}
	var s model.Space
	query.First(&s)
	return s.Id, nil
}
