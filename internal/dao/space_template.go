package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
	"fmt"
)

func FindAllTemplate() (space_template []*model.SpaceTemplate) {
	db.DB.Where("status = 1").Find(&space_template)
	return
}

func FindOneTemplateByName(name string) (space_template model.SpaceTemplate, ok bool) {
	fmt.Println(name, "8888888", db.DB)
	db.DB.Where("name = ? AND status = 1", name).First(&space_template)
	return space_template, space_template == model.SpaceTemplate{}
}

func FindOneTemplateById(st_id string) (space_template model.SpaceTemplate, ok bool) {
	db.DB.Where("id = ? AND status = 1", st_id).First(&space_template)
	return space_template, space_template == model.SpaceTemplate{}
}

func InsertSpaceTemplate(space_template *model.SpaceTemplate) (st_id uint32, err error) {
	query := db.DB.Create(space_template)
	err = query.Error
	if err != nil {
		return 0, err
	}
	var st model.SpaceTemplate
	query.Where("name = ?", space_template.Name).First(&st)
	return st.Id, nil
}

func UpdateSpaceTemplate(req model.SpaceTemplateCreateOption, st_id string) (err error) {
	var space_template model.SpaceTemplate
	if err := db.DB.Model(&space_template).Where("id = ?", st_id).Updates(map[string]interface{}{
		"kind_id":     req.KindId,
		"name":        req.Name,
		"description": req.Description,
		"tags":        req.Tags,
		"image":       req.Image,
		"avatar":      req.Avatar,
		"status":      1,
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSpaceTemplate(st_id string) (err error) {
	var space_template model.SpaceTemplate
	if err := db.DB.Model(&space_template).Where("id = ?", st_id).Update("status", 0).Error; err != nil {
		return err
	}
	return nil
}
