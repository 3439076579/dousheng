package controller

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"awesomeProject/model/video/response"
	"awesomeProject/service"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetVideoList(router *gin.Context) {

	var VideoList []video.Video
	var err error
	var HasToken bool
	var UserID int64

	id, ok := router.Get("user_id")
	if !ok {
		HasToken = false
	} else {
		HasToken = true
		UserID = id.(int64)
	}

	VideoList, err = service.GetVideoFeedService(HasToken, UserID)
	if err != nil {
		router.JSON(http.StatusOK, &response.ResponseGetVideo{
			Response: common.Response{
				StatusCode: -1,
				StatusMsg:  "返回信息错误",
			},
			NextTime:  time.Now().Unix(),
			VideoList: nil,
		})
		return
	}

	router.JSON(http.StatusOK, &response.ResponseGetVideo{
		VideoList: VideoList,
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		NextTime: time.Now().Unix(),
	})
	return
}

// GetPublishedVideo 获取发布过的视频，登录用户才可以获取发布过的视频
func GetPublishedVideo(router *gin.Context) {

	var user_ user.User
	// 从gin.Context中获取到user_id，根据user_id查找作品列表
	tmp, ok := router.Get("user_id")
	if !ok {
		fmt.Println("未接受到user_id")
		router.JSON(http.StatusOK, &response.ResponseGetPublishedVideo{
			Response:  common.Response{StatusCode: -1, StatusMsg: "未登录用户不可使用此功能"},
			VideoList: nil,
		})
		return
	}
	// 获取到了传入的user_id参数
	user_.UserId, _ = tmp.(int64)

	videoList, err := service.GetPublishedVideoService(&user_)
	if err != nil {
		fmt.Println("查找视频失败")
		return
	}

	router.JSON(http.StatusOK, &response.ResponseGetPublishedVideo{
		VideoList: videoList,
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
	})
	return
}

func PublishVideo(router *gin.Context) {

	var video_ video.Video

	video_.Title = router.PostForm("title")
	id, ok := router.Get("user_id")
	if !ok {
		router.JSON(http.StatusOK, &common.Response{
			StatusCode: -1,
			StatusMsg:  "token is incorrect or token loss",
		})
		return
	}
	video_.AuthorId = id.(int64)
	// 获取视频资源
	videoFile, err := router.FormFile("data")
	if err != nil {
		router.JSON(http.StatusOK, &common.Response{
			StatusCode: -1,
			StatusMsg:  "fail to publish,lack of video",
		})
		fmt.Println(err)
		return
	}

	// 调用UpLoadFile函数进行上传文件，返回播放的URL
	loadFile, err := utils.UpLoadFile(video_.Title, *videoFile, videoFile.Size)
	if err != nil {
		fmt.Println(err, loadFile)
		return
	}

	err = service.PublishVideoService(&video_)
	if err != nil {
		router.JSON(http.StatusOK, &common.Response{
			StatusCode: -1,
			StatusMsg:  "fail to publish",
		})
		return
	}

	// 未出错返回"success to upload file"
	router.JSON(http.StatusOK, &common.Response{
		StatusCode: 0,
		StatusMsg:  "success to upload file",
	})
}
