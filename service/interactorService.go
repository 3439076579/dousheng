package service

import (
	"awesomeProject/dao"
	"awesomeProject/model/interactor"
	"awesomeProject/model/interactor/request"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"awesomeProject/redis_"
	"awesomeProject/utils"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"strings"
)

/*
step:
for favourite:
	1.查看是否已经点过赞了
	2.没点过就使用Redis加入到set中
for cancel favourite:
	1.查看是否没有点过赞
 	2.如果点过了，就从Redis去除这个记录
*/
// 利用redis的map结构来存储：

const (
	// MAP_KEY_VIDEO_LIKED key:MAP_VIDEO_LIKED field:UserID::VideoID value:0或1
	MAP_KEY_VIDEO_LIKED = "MAP_VIDEO_LIKED"
)

// FavouriteActionService 用于进行点赞操作，点赞操作是并发安全的
func FavouriteActionService(favourite *request.RequestFavourite) error {

	// 转化结构为FavouriteRelation
	var FavouriteModel interactor.FavouriteRelation

	FavouriteModel.VideoID = favourite.VideoID
	FavouriteModel.UserID = favourite.UserLoginID

	switch favourite.ActionType {
	case "1":
		return AddFavorite(&FavouriteModel)

	case "2":
		return CancelFavorite(&FavouriteModel)
	default:
		return utils.InvalidFavoriteAction
	}

}

func AddFavorite(favoriteModel *interactor.FavouriteRelation) error {

	// 创建Redis服务
	var r = redisService.RedisService{
		Ctx: context.Background(),
	}

	result := redisService.GlobalRedis.HGet(r.Ctx, MAP_KEY_VIDEO_LIKED,
		strconv.FormatInt(favoriteModel.UserID, 10)+
			"::"+strconv.FormatInt(favoriteModel.VideoID, 10))

	if result.Val() == "1" {
		return utils.DuplicateAddFavorite
	}

	redisService.GlobalRedis.HSet(r.Ctx, MAP_KEY_VIDEO_LIKED, strconv.FormatInt(favoriteModel.UserID, 10)+
		"::"+strconv.FormatInt(favoriteModel.VideoID, 10), 1)
	redisService.GlobalRedis.IncrBy(r.Ctx,
		strconv.FormatInt(favoriteModel.UserID, 10)+"::"+
			"USER_FAVORITE_COUNT", 1)
	lock := redisService.GetDistributedLock(30)
	lock.Lock()
	redisService.GlobalRedis.IncrBy(r.Ctx, strconv.FormatInt(favoriteModel.VideoID, 10)+"::"+
		"VIDEO_FAVORITE_COUNT", 1)
	err := lock.Unlock()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CancelFavorite(favoriteModel *interactor.FavouriteRelation) error {

	// 创建Redis服务
	var r = redisService.RedisService{
		Ctx: context.Background(),
	}

	result := redisService.GlobalRedis.HGet(r.Ctx, MAP_KEY_VIDEO_LIKED,
		strconv.FormatInt(favoriteModel.UserID, 10)+
			"::"+strconv.FormatInt(favoriteModel.VideoID, 10))

	if result.Val() == "1" {
		redisService.GlobalRedis.HDel(r.Ctx, MAP_KEY_VIDEO_LIKED,
			strconv.FormatInt(favoriteModel.UserID, 10)+
				"::"+strconv.FormatInt(favoriteModel.VideoID, 10))
		redisService.GlobalRedis.IncrBy(r.Ctx,
			strconv.FormatInt(favoriteModel.UserID, 10)+"::"+
				"USER_FAVORITE_COUNT", -1)
		lock := redisService.GetDistributedLock(30)
		lock.Lock()
		redisService.GlobalRedis.IncrBy(r.Ctx, strconv.FormatInt(favoriteModel.VideoID, 10)+"::"+
			"VIDEO_FAVORITE_COUNT", -1)
		err := lock.Unlock()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	return utils.CancelFavoriteErr

}

func GetFavouriteListService(req *request.RequestGetFavouriteList) ([]video.Video, error) {

	// 获取到用户ID
	UserID := req.UserLoginID

	// 根据用户ID获取Redis数据
	redisStore := &redisService.RedisService{
		Ctx: context.Background(),
	}

	stringMap := redisService.GlobalRedis.HGetAll(redisStore.Ctx, MAP_KEY_VIDEO_LIKED)
	result, err := stringMap.Result()
	if err != nil {
		return nil, err
	}

	var VideoIDArray []string

	for k, v := range result {
		if v == "1" &&
			strconv.FormatInt(UserID, 10) == strings.Split(k, "::")[0] {
			VideoIDArray = append(VideoIDArray, strings.Split(k, "::")[1])
		}
	}

	// 根据resultSlice获取Videos
	videoDao := dao.NewVideoDao(true)
	defer videoDao.DB.Commit()

	videos, err := videoDao.SearchVideoByFavouriteID(VideoIDArray)
	if err != nil {
		videoDao.DB.Rollback()
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		videos[i].IsFavorite = true
	}

	return videos, nil

}

func CommentService(comment *request.RequestPostComment) (interactor.Comment, error) {
	var CommentModel interactor.Comment
	var err error

	CommentModel.VideoID = comment.VideoID
	CommentModel.UserID = comment.UserID

	switch comment.ActionType {
	case "1":
		CommentModel.Content = comment.Comment
		err = PostComment(&CommentModel)

	case "2":
		CommentModel.ID = uint(comment.CommentID)
		err = DeleteComment(&CommentModel)
	}

	if err != nil {
		return interactor.Comment{}, err
	}

	return CommentModel, nil
}

// DeleteComment 用于执行删除评论的操作
func DeleteComment(CommentModel *interactor.Comment) error {
	interactDao := dao.NewInteractDao(true)
	defer interactDao.DB.Commit()

	err := interactDao.DeleteCommentByID(CommentModel)
	if err != nil {
		interactDao.DB.Rollback()
		return err
	}
	videoDao := dao.NewVideoDaoFromDB(interactDao.DB)
	err = videoDao.UpdateCommentCount(-1, CommentModel.VideoID)
	if err != nil {
		interactDao.DB.Rollback()
		return err
	}

	r := redisService.RedisService{
		PreKey: strconv.FormatInt(CommentModel.VideoID, 10) + ":comment",
		Ctx:    context.Background(),
	}

	redisService.GlobalRedis.SRem(r.Ctx, r.PreKey, strconv.FormatInt(CommentModel.UserID, 10)+
		","+CommentModel.CreatedAt.Format("2006-01-02 15:04:05")+
		","+CommentModel.Content)

	return nil

}

// PostComment 用于执行增加评论
func PostComment(CommentModel *interactor.Comment) error {

	CommentModel.ID = uint(utils.GenerateUuid())

	// 将评论插入数据库
	interactDao := dao.NewInteractDao(true)

	err := interactDao.CreateComment(CommentModel)
	if err != nil {
		interactDao.DB.Rollback()
		return err
	}
	videoDao := dao.NewVideoDaoFromDB(interactDao.DB)
	err = videoDao.UpdateCommentCount(1, CommentModel.VideoID)
	if err != nil {
		interactDao.DB.Rollback()
		return err
	}
	interactDao.DB.Commit()
	// 把评论缓存到Redis
	r := redisService.RedisService{
		PreKey: strconv.FormatInt(CommentModel.VideoID, 10) + ":comment",
		Ctx:    context.Background(),
	}

	redisService.GlobalRedis.SAdd(r.Ctx, r.PreKey, strconv.FormatInt(CommentModel.UserID, 10)+
		","+CommentModel.CreatedAt.Format("2006-01-02 15:04:05")+
		","+CommentModel.Content)

	return nil

}

func GetCommentListService(comment *request.RequestGetComment) ([]interactor.Comment, error) {
	var CommentList []interactor.Comment

	// 尝试去Redis中获取
	r := redisService.RedisService{
		Ctx:    context.Background(),
		PreKey: strconv.FormatInt(comment.VideoID, 10) + ":comment",
	}

	RedisResult := r.SGet(r.PreKey)
	// Redis中获取到评论
	if RedisResult.Err() != redis.Nil && len(RedisResult.Val()) != 0 {
		CommentList = utils.ParseStringSliceToCommentList(RedisResult.Val())
		// 从数据库中获取用户信息
		userDao := dao.NewUserDao(true)
		defer userDao.DB.Commit()
		for i := 0; i < len(CommentList); i++ {
			var userInfo *user.User = &user.User{UserId: CommentList[i].UserID}
			err := userDao.SearchUserInfoByUserId(userInfo)
			if err != nil {
				userDao.DB.Rollback()
				return CommentList, err
			}
			CommentList[i].User = *userInfo
		}
		return CommentList, nil
	}
	// 从Redis中获取不到评论，从数据库中获取
	interactDao := dao.NewInteractDao(true)
	defer interactDao.DB.Commit()

	err := interactDao.SearchCommentByVideoIDInBatch(&CommentList, comment.VideoID)
	if err != nil {
		interactDao.DB.Rollback()
		return CommentList, err
	}

	userDao := dao.NewUserDaoFromDB(interactDao.DB)
	for i := 0; i < len(CommentList); i++ {
		var userInfo *user.User = &user.User{UserId: CommentList[i].UserID}
		err := userDao.SearchUserInfoByUserId(userInfo)
		if err != nil {
			userDao.DB.Rollback()
			return CommentList, err
		}
		CommentList[i].User = *userInfo
		// 把评论缓存到Redis中
		/* key的组装
		UserID,CreatedAt,Content
		*/

		r.SAdd(r.PreKey, strconv.FormatInt(CommentList[i].UserID, 10)+
			","+CommentList[i].CreatedAt.Format("2006-01-02 15:04:05")+
			","+CommentList[i].Content)
	}

	return CommentList, nil

}
