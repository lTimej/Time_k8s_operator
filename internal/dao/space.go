package dao

import (
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
	"time"
)

const (
	TemplateDeleted = iota
	TemplateUsing
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

func FindAllSpaceByUserId(user_id uint32) (spaces []*model.Space) {
	db.DB.Where("user_id = ? AND status != 0", user_id).Find(&spaces)
	return
}

func FindOneByUserIdAndName(user_id uint32, name string) bool {
	var space model.Space
	db.DB.Where("user_id = ? AND name = ? AND status != ?", user_id, name, model.SpaceStatusDeleted).First(&space)
	return model.Space{} == space
}

func FindCountByUserId(user_id uint32) int64 {
	var count int64
	db.DB.Table("space").Where("user_id = ? AND status != 0",user_id).Count(&count)
	return count
}

func InsertSpace(space *model.Space) (space_id uint32, err error) {
	query := db.DB.Create(space)
	err = query.Error
	if err != nil {
		return 0, err
	}
	var s model.Space
	query.Where("sid = ? AND name = ?", space.Sid, space.Name).First(&s)
	return s.Id, nil
}

func UpdateSpaceStatusAndRunningStatus(space_id uint32, status int, running_status uint32) error {
	var space model.Space
	if err := db.DB.Model(&space).Where("id = ?", space_id).Updates(map[string]interface{}{"status": status, "running_status": running_status}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSpaceStatus(space_id uint32, status int) error {
	var space model.Space
	if err := db.DB.Model(&space).Where("id = ?", space_id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSpaceRunningStatus(space_id uint32, running_status int) error {
	space := FindSpaceOneById(space_id)
	now := time.Now()
	space.StopTime = now
	space.TotalTime = now.Sub(space.CreateTime)
	space.RunningStatus = uint32(running_status)
	if err := db.DB.Save(&space).Error; err != nil {
		return err
	}
	// if err := db.DB.Model(&space).Where("id = ?", space_id).Updates(map[string]interface{}{"running_status": running_status, "stop_time": now}).Error; err != nil {
	// 	return err
	// }
	return nil
}

func FindSpaceOneByIdAndUserId(id, user_id uint32) (space model.Space, ok bool) {
	db.DB.Where("id = ? AND user_id = ?", id, user_id).First(&space)
	return space, model.Space{} == space
}

func FindSpaceOneById(id uint32) (space model.Space) {
	db.DB.Where("id = ?", id).First(&space)
	return space
}

func DeleteSpaceById(space_id uint32) error {
	var space model.Space
	if err := db.DB.Model(&space).Where("id = ?", space_id).Update("status", model.SpaceStatusDeleted).Error; err != nil {
		return err
	}
	return nil
}
