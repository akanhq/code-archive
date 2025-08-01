package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PostID    primitive.ObjectID `bson:"post_id"`                     //文章id
	AuthorID  primitive.ObjectID `bson:"author_id"`                   //评论人id
	Author    string             `bson:"author"`                      //评论人
	Content   string             `bson:"content" validate:"required"` //内容
	CreatedAt time.Time          `bson:"created_at"`                  //创建时间
}
