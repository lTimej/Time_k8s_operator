package main

import (
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/model"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("data_table.yaml")
	viper.AddConfigPath("../../templates")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("file not found.")
			return
		}
		panic(err)
	}
	templates := viper.GetStringMap("space_template")
	for _, val := range templates {
		template := val.(map[string]interface{})
		fmt.Println(template["avatar"])
		now := time.Now()
		space_template := &model.SpaceTemplate{
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
		fmt.Println(space_template.Name)
		_, ok := dao.FindOneTemplateByName(space_template.Name)
		if !ok {
			fmt.Println("模板名称已存在")
			return
		}
		fmt.Println("111111")
		st_id, err := dao.InsertSpaceTemplate(space_template)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(st_id)
	}
}
