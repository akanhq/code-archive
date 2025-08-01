package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`                   //ID
	Username  string             `bson:"username" validate:"required"`    //用户名称
	Password  string             `bson:"password" validate:"required"`    //密码
	Email     string             `bson:"email" validate:"required,email"` //邮箱
	CreatedAt time.Time          `bson:"created_at"`                      //创建时间
	UpdatedAt time.Time          `bson:"updated_at"`                      //更新时间
}
