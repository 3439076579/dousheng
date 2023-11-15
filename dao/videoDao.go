package dao

import (
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"awesomeProject/utils"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type VideoDao struct {
	DB *gorm.DB
}

func NewVideoDao(UserTransaction bool) VideoDao {
	return VideoDao{DB: utils.GetDB(UserTransaction)}
}

// NewVideoDaoFromDB 复用DB连接的函数
func NewVideoDaoFromDB(db *gorm.DB) VideoDao {
	return VideoDao{DB: db}
}

// SearchPublishedListByAuthorId 通过作者ID获取该用户发布过的所有视频
func (v VideoDao) SearchPublishedListByAuthorId(user_ *user.User) ([]video.Video, error) {

	var videolist []video.Video

	res := v.DB.Model(&video.Video{}).
		Preload("Author").
		Where("author_id=?", user_.ID).
		Find(&videolist)
	if res.Error != nil {
		return nil, res.Error
	} else {
		if res.RowsAffected == 0 {
			return nil, utils.RecordNotFound
		}
	}
	return videolist, nil
}

// SearchVideoByRandom 随机选择8个视频返回
func (v VideoDao) SearchVideoByRandom() ([]video.Video, error) {

	var videolist []video.Video
	// 需要引入随机性
	res := v.DB.Model(&video.Video{}).
		Preload("Author").
		Find(&videolist)
	if res.Error != nil {
		return nil, res.Error
	} else {
		if res.RowsAffected == 0 {
			return nil, utils.RecordNotFound
		}
	}

	return videolist, nil

}

// InsertVideo 发布视频
func (v VideoDao) InsertVideo(video_ *video.Video) error {

	res := v.DB.Model(&video.Video{}).
		Preload("Author").
		Create(&video_)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (v VideoDao) SearchVideoByID(videoID int64) (video.Video, error) {

	var Video video.Video

	res := v.DB.Model(&video.Video{}).
		Preload("Author").
		Where("id=?", videoID).
		Find(&Video)
	if res.Error != nil {
		return Video, res.Error
	} else {
		if res.RowsAffected == 0 {
			return Video, utils.RecordNotFound
		}
	}
	return Video, nil
}

// SearchVideoByFavouriteID 通过点赞视频的ID查找喜欢列表
func (v VideoDao) SearchVideoByFavouriteID(favouriteID []string) ([]video.Video, error) {
	var VideoList []video.Video

	for i := 0; i < len(favouriteID); i++ {

		id, err := strconv.ParseInt(favouriteID[i], 10, 64)
		if err != nil {
			log.Println("id converting occurs error:", err)
			return VideoList, err
		}
		Video, err := v.SearchVideoByID(id)
		if err != nil {
			log.Println("Search Video occurs error:", err)
			return VideoList, err
		}
		VideoList = append(VideoList, Video)
	}

	return VideoList, nil

}

func (v VideoDao) UpdateCommentCountAndFavouriteCount(videoModel video.Video) error {
	res := v.DB.Model(&video.Video{}).
		Where("id=?", videoModel.ID).
		Updates(map[string]interface{}{"comment_count": videoModel.CommentCount,
			"favorite_count": videoModel.FavoriteCount})
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil

}

func (v VideoDao) UpdateCommentCount(count int64, VideoID int64) error {
	res := v.DB.Model(video.Video{}).
		Where("id=?", VideoID).
		Update("comment_count", gorm.Expr("comment_count+?", count))
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil
}

func (v VideoDao) UpdateFavoriteCount(favoriteCount int64, VideoID int64) error {
	res := v.DB.Model(video.Video{}).
		Where("id=?", VideoID).
		Update("favorite_count", favoriteCount)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil

}
