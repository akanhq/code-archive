package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`               //ID
	Title     string             `bson:"title" validate:"required"`   //标题
	Content   string             `bson:"content" validate:"required"` //内容
	AuthorID  primitive.ObjectID `bson:"author_id"`                   //创建
	Author    string             `bson:"author"`                      //创建
	Views     int                `bson:"views"`                       //点击量
	CreatedAt time.Time          `bson:"created_at"`                  //创建时间
	UpdatedAt time.Time          `bson:"updated_at"`                  //更新时间
}
