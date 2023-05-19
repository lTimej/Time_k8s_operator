package model

import "time"

type User struct {
	Id         uint32    `gorm:"type:bigint(20)"`
	Uid        string    `gorm:"type:varchar(30);not null;unique"`
	Username   string    `gorm:"type:varchar(255);not null;unique"`
	Password   string    `gorm:"type:char(255);not null"`
	Nickname   string    `gorm:"type:varchar(255);"`
	Email      string    `gorm:"type:varchar(255);default:NULL"`
	Phone      string    `gorm:"type:char(11)"`
	Avatar     string    `gorm:"type:varchar(255);"`
	CreateTime time.Time `gorm:"type:datetime; DEFAULT CURRENT_TIMESTAMP"`
	DeleteTime time.Time `gorm:"type:datetime; DEFAULT NULL"`
	Status     uint32    `gorm:"type:tinyint(1)"` // 状态 0正常 1已删除

	Token string `gorm:"type:varchar(30)"`
}

func (m *User) TableName() string {
	return "user"
}

type RegisterInfo struct {
	Nickname   string `json:"nickname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
	Email      string `json:"email"`
}
