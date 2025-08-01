package models

import "time"

type CommonField struct {
	Id         uint64    `json:"id" gorm:"primarykey"`              // 主键id
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"` // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"` // 更新时间
	CreateId   uint64    `json:"create_id"`                         // 创建人id
	CreateName string    `json:"create_name"`                       // 创建人真实姓名
	UpdateId   uint64    `json:"update_id"`                         // 更新人Id',
	UpdateName string    `json:"update_name"`                       // 更新人真实姓名
	DelFlag    int32     `json:"del_flag"`                          // 删除状态 1：未删除 -1：已删除
}
