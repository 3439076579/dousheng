package timer

import (
	"awesomeProject/dao"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"awesomeProject/redis_"
	"context"
	"strconv"
	"strings"
	"time"
)

// InitTimer 初始化所有的定时器
func InitTimer() {

	timer := time.NewTimer(time.Second * 1)

	for {
		select {
		case <-timer.C:
			ReWriteToDB()
			timer.Reset(time.Second * 1)
		default:
			time.Sleep(time.Millisecond * 999)
		}
	}

}

/*
	二分查找：[left,right)

*/

func ExistsVideoInVideoList(videoList []video.Video, VideoId int64) int64 {

	var left int = 0
	var right int = len(videoList)

	for left < right {
		mid := left + (right-left)>>1

		if videoList[mid].ID == VideoId {
			return int64(mid)
		} else if videoList[mid].ID > VideoId {
			right = mid
		} else {
			left = mid + 1
		}

	}
	return -1

}

/*

	需要回写进数据库的field：
			- VideoID:VideoID  ------>代表当前视频被喜欢的人数 ----->video.FavoriteCount
			- UserLoginID:VideoID ----->代表该用户喜欢哪些视频 ----->user.FavoriteCount

*/

func ReWriteToDB() {

	var cursor uint64

	r := redisService.RedisService{
		Ctx: context.Background(),
	}

	// step1:扫描出所有*::USER_FAVORITE_COUNT
	result, cursor, err := redisService.GlobalRedis.
		Scan(r.Ctx, cursor, "*::USER_FAVORITE_COUNT", 1).Result()
	if err != nil {
		return
	}

	userDao := dao.NewUserDao(true)
	defer userDao.DB.Commit()

	for i := 0; i < len(result); i++ {
		UserID, _ := strconv.ParseInt(strings.Split(result[i], "::")[0], 10, 64)
		var userModel user.User
		userModel.UserId = UserID
		userModel.FavoriteCount, _ = strconv.ParseInt(
			redisService.GlobalRedis.Get(r.Ctx, result[i]).Val(), 10, 64)

		err := userDao.IncreaseFavoriteCount(userModel)
		if err != nil {
			userDao.DB.Rollback()
			return
		}
	}
	result, cursor, err = redisService.GlobalRedis.
		Scan(r.Ctx, cursor, "*::VIDEO_FAVORITE_COUNT", 1).Result()
	if err != nil {
		return
	}

	videoDao := dao.NewVideoDaoFromDB(userDao.DB)

	for i := 0; i < len(result); i++ {
		VideoID, _ := strconv.ParseInt(strings.Split(result[i], "::")[0], 10, 64)
		var videoModel video.Video
		videoModel.ID = VideoID
		videoModel.FavoriteCount, _ = strconv.ParseInt(
			redisService.GlobalRedis.Get(r.Ctx, result[i]).Val(), 10, 64)
		err := videoDao.UpdateFavoriteCount(videoModel.FavoriteCount, videoModel.ID)
		if err != nil {
			userDao.DB.Rollback()
			return
		}
	}

}
