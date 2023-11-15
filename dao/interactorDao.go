package dao

import (
	"awesomeProject/model/interactor"
	"awesomeProject/utils"
	"gorm.io/gorm"
)

func NewInteractDao(UseTransaction bool) InteractDao {
	return InteractDao{DB: utils.GetDB(UseTransaction)}
}

func NewInteractDaoFromDB(db *gorm.DB) InteractDao {
	return InteractDao{DB: db}
}

type InteractDao struct {
	DB *gorm.DB
}

func (i InteractDao) ExistFavourite(Favourite *interactor.FavouriteRelation) (bool, error) {

	var NewFavourite *interactor.FavouriteRelation

	if res := i.DB.Model(&Favourite).
		Where("video_id=? AND user_id=?", Favourite.VideoID, Favourite.UserID).
		Find(&NewFavourite); res.Error != nil {
		return false, res.Error
	} else {
		if res.RowsAffected == 0 {
			return false, nil
		} else {
			return true, nil
		}
	}
}

func (i InteractDao) InsertFavouriteRelation(Favourite *interactor.FavouriteRelation) error {
	res := i.DB.Model(&Favourite).
		Create(&Favourite)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (i InteractDao) CreateComment(CommentModel *interactor.Comment) error {
	res := i.DB.Model(&CommentModel).
		Create(&CommentModel)

	if res.Error != nil {
		return res.Error
	}

	return nil

}

func (i InteractDao) DeleteCommentByID(CommentModel *interactor.Comment) error {
	res := i.DB.Model(&CommentModel).
		Where("id=?", CommentModel.ID).
		Delete(&CommentModel)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil
}

func (i InteractDao) SearchCommentByVideoID(CommentModel *interactor.Comment) error {

	res := i.DB.Model(&CommentModel).
		Where("video_id=?", CommentModel.VideoID).
		Find(&CommentModel)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
		return nil
	}

}

func (i InteractDao) SearchCommentByVideoIDInBatch(CommentList *[]interactor.Comment, videoId int64) error {

	rows, err := i.DB.Model(&interactor.Comment{}).Preload("User").
		Rows()
	defer rows.Close()

	for rows.Next() {
		var CommentModel interactor.Comment

		i.DB.ScanRows(rows, &CommentModel)

		*CommentList = append(*CommentList, CommentModel)

	}
	if err != nil {
		return err
	}
	return nil

}
