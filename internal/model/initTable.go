package model

import (
	"Time_k8s_operator/internal/dao/db"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func InitTable() error {
	// if err := db.DB.AutoMigrate(&User{}, &TemplateKind{}, &SpaceTemplate{}, &SpaceSpec{}, &Space{}); err != nil {
	// 	return err
	// }
	if err := db.DB.AutoMigrate(&Space{}); err != nil {
		return err
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
