package interactor

import (
	"awesomeProject/model/user"
	"gorm.io/gorm"
)

type FavouriteRelation struct {
	UserID  int64 `gorm:"column:user_id"`
	VideoID int64 `gorm:"column:video_id"`
}

func (FavouriteRelation) TableName() string {
	return "favourite_video_user"
}

type Comment struct {
	gorm.Model
	Content string    `gorm:"column:content"`
	User    user.User `gorm:"foreignKey:UserID;references:ID"`
	UserID  int64     `gorm:"column:user_id"`
	VideoID int64     `gorm:"column:video_id"`
}
