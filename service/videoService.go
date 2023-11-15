package service

import (
	"awesomeProject/dao"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"awesomeProject/redis_"
	"context"
	"strconv"
	"strings"
)

func GetPublishedVideoService(user_ *user.User) ([]video.Video, error) {

	var VideoList []video.Video

	// 查询数据库
	userDao := dao.NewUserDao(true)
	defer userDao.DB.Commit()
	// 从数据库中找到user信息
	err := userDao.SearchUserInfoByUserId(user_)
	if err != nil {
		userDao.DB.Rollback()
		return nil, err
	}

	// 复用userDao中的db
	videoDao := dao.NewVideoDaoFromDB(userDao.DB)
	// 根据AuthorId查找到该用户发布的视频列表
	VideoList, err = videoDao.SearchPublishedListByAuthorId(user_)
	if err != nil {
		videoDao.DB.Rollback()
		return nil, err
	}
	r := redisService.RedisService{
		Ctx: context.Background(),
	}

	slice, err := redisService.GlobalRedis.HGetAll(r.Ctx, MAP_KEY_VIDEO_LIKED).Result()
	if err != nil {
		return nil, err
	}
	for k, v := range slice {
		if v == "1" &&
			strings.Split(k, "::")[0] == strconv.FormatInt(user_.UserId, 10) {
			for i := 0; i < len(VideoList); i++ {
				if strconv.FormatInt(VideoList[i].ID, 10) == strings.Split(k, "::")[0] {
					VideoList[i].IsFavorite = true
				}
			}
		}
	}

	return VideoList, nil
}

func GetVideoFeedService(HasToken bool, UserID ...int64) ([]video.Video, error) {

	var VideoList []video.Video

	videoDao := dao.NewVideoDao(true)
	defer videoDao.DB.Commit()
	VideoList, err := videoDao.SearchVideoByRandom()
	if err != nil {
		videoDao.DB.Rollback()
		return nil, err
	}
	r := redisService.RedisService{
		Ctx: context.Background(),
	}

	// 如果有Token，代表用户有ID
	if HasToken {

		slice, err := redisService.GlobalRedis.HGetAll(r.Ctx, MAP_KEY_VIDEO_LIKED).Result()
		if err != nil {
			return nil, err
		}
		for k, v := range slice {
			if v == "1" &&
				strings.Split(k, "::")[0] == strconv.FormatInt(UserID[0], 10) {
				for i := 0; i < len(VideoList); i++ {
					if strconv.FormatInt(VideoList[i].ID, 10) == strings.Split(k, "::")[0] {
						VideoList[i].IsFavorite = true
					}
				}
			}
		}

	}
	return VideoList, nil
}

func PublishVideoService(video_ *video.Video) error {

	var userModel user.User

	userModel.ID = video_.AuthorId

	// 执行插入逻辑即可
	videoDao := dao.NewVideoDao(true)
	defer videoDao.DB.Commit()

	err := videoDao.InsertVideo(video_)
	if err != nil {
		videoDao.DB.Rollback()
		return err
	}
	userDao := dao.NewUserDaoFromDB(videoDao.DB)

	err = userDao.IncreaseWorkCount(userModel, 1)
	if err != nil {
		videoDao.DB.Rollback()
		return err
	}

	return nil
}
