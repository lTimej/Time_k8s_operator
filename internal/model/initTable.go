package model

import (
	"Time_k8s_operator/internal/dao/db"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func InitTable() error {
	if err := db.DB.AutoMigrate(&User{}, &TemplateKind{}, &SpaceTemplate{}, &SpaceSpec{}, &Space{}); err != nil {
		return err
	}
	// if err := db.DB.AutoMigrate(&Space{}); err != nil {
	// 	return err
	// }
	return nil
}

func initSpaceTemplate() error {
	templates := viper.GetStringMap("space_template")
	for _, val := range templates {
		template := val.(map[string]interface{})
		now := time.Now()
		space_template := &SpaceTemplate{
			Name:        template["name"].(string),
			KindId:      uint32(template["kind_id"].(int)),
			Description: template["description"].(string),
			Tags:        template["tags"].(string),
			Image:       template["image"].(string),
			Status:      1,
			Avatar:      template["avatar"].(string),
			CreateTime:  now,
			DeleteTime:  now,
		}
		var stt SpaceTemplate
		db.DB.Where("name = ? AND status = 1", space_template.Name).First(&stt)
		if stt.Name != "" {
			fmt.Println("模板名称已存在", stt.Name)
			continue
		}
		if err := db.DB.Create(space_template).Error; err != nil {
			fmt.Println(err)
			return err
		}

	}
	return nil
}

func initTemplateKind() error {
	kinds := viper.GetStringMap("template_kin")
	for _, val := range kinds {
		kind := val.(map[string]interface{})
		template_kind := &TemplateKind{
			Name: kind["name"].(string),
		}
		var tkk TemplateKind
		db.DB.Where("name = ?", template_kind.Name).First(&tkk)
		if tkk.Name != "" {
			fmt.Println("模板类型已存在", tkk.Name)
			continue
		}
		if err := db.DB.Create(template_kind).Error; err != nil {
			fmt.Println(err)
			return err
		}

	}
	return nil
}

func initSpaceSpec() error {
	specs := viper.GetStringMap("space_spec")
	for _, val := range specs {
		spec := val.(map[string]interface{})
		space_spec := &SpaceSpec{
			Name:        spec["name"].(string),
			CpuSpec:     spec["cpu_spec"].(string),
			MemSpec:     spec["mem_spec"].(string),
			StorageSpec: spec["storage_spec"].(string),
			Description: spec["description"].(string),
		}
		if err := db.DB.Create(space_spec).Error; err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
func InitTableData() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("data_table.yaml")
	viper.AddConfigPath("./templates")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("file not found.")
			return errors.New("数据库初始化配置文件不存在")
		}
		panic(err)
	}
	err := initSpaceTemplate()
	if err != nil {
		return err
	}
	err = initTemplateKind()
	if err != nil {
		return err
	}
	err = initSpaceSpec()
	if err != nil {
		return err
	}
	return nil
}
