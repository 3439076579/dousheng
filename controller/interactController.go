package controller

import (
	"awesomeProject/model/common"
	"awesomeProject/model/interactor/request"
	"awesomeProject/model/video/response"
	"awesomeProject/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Favourite(ctx *gin.Context) {
	// 获取相应的参数
	var favouriteReq request.RequestFavourite
	favouriteReq.ActionType = ctx.Query("action_type")
	id, err := strconv.Atoi(ctx.Query("video_id"))
	if err != nil {
		log.Println("videoID is incorrect")
		return
	}
	favouriteReq.VideoID = int64(id)
	UserID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("lack of user_id in token")
		return
	}
	favouriteReq.UserLoginID = UserID.(int64)

	// 调用service层完成业务逻辑
	err = service.FavouriteActionService(&favouriteReq)
	if err != nil {
		//点赞失败
		log.Println("点赞失败", err)
		ctx.JSON(http.StatusOK, &common.Response{
			StatusCode: -1,
			StatusMsg:  "fail to favourite",
		})
		return
	}

	ctx.JSON(http.StatusOK, &common.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})

}

func GetFavouriteList(ctx *gin.Context) {

	var getFav request.RequestGetFavouriteList
	// 获取到UserID
	ID, ok := ctx.Get("user_id")
	if !ok {
		log.Println("lack of user_id in token")
		return
	}
	getFav.UserLoginID = ID.(int64)

	// 调用Service层
	videos, err := service.GetFavouriteListService(&getFav)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, &response.ResponseGetPublishedVideo{
			Response:  common.Response{StatusCode: -1, StatusMsg: "出现错误"},
			VideoList: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, &response.ResponseGetPublishedVideo{
		Response:  common.Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videos,
	})
}

// CommentHandler 用于增加评论和删除评论的Handler
func CommentHandler(ctx *gin.Context) {

	var PostCommentReq request.RequestPostComment
	// 获取参数UserLoginID
	id, ok := ctx.Get("user_id")
	if !ok {
		log.Println("lack of user id in token ")
		return
	}
	PostCommentReq.UserID = id.(int64)
	// 获取参数VideoID
	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("occurs error", err)
		return
	}
	PostCommentReq.VideoID = videoId
	// 根据ActionType选择性获取对应的参数
	PostCommentReq.ActionType = ctx.Query("action_type")
	if PostCommentReq.ActionType == "1" {
		PostCommentReq.Comment = ctx.Query("comment_text")
	}
	if PostCommentReq.ActionType == "2" {
		commentId, err := strconv.ParseInt(ctx.Query("comment_id"), 10, 64)
		if err != nil {
			log.Println("occurs error", err)
			return
		}
		PostCommentReq.CommentID = commentId
	}

	// 调用Service层进行业务处理
	CommentModel, err := service.CommentService(&PostCommentReq)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, &response.ResponsePostComment{
			Response: common.Response{StatusCode: -1, StatusMsg: err.Error()},
			Comment:  nil,
		})
		return
	}
	// CommentModel转化为ResponseComment
	var responseComment response.ResponseComment

	responseComment.CommentID = int64(CommentModel.ID)
	responseComment.User = CommentModel.User
	responseComment.Content = CommentModel.Content
	responseComment.Created = CommentModel.CreatedAt.Format("2006-01-02 15:04:05")

	ctx.JSON(http.StatusOK, &response.ResponsePostComment{
		Response: common.Response{StatusCode: 0, StatusMsg: "success"},
		Comment:  &responseComment,
	})

}

func GetCommentList(ctx *gin.Context) {
	var getCommentReq request.RequestGetComment

	videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 64)
	if err != nil {
		log.Println("illegal video_id")
	}

	getCommentReq.VideoID = videoId

	id, ok := ctx.Get("user_id")
	if !ok {
		log.Println("lack of user id in token")
		return
	}

	getCommentReq.UserID = id.(int64)

	commentList, err := service.GetCommentListService(&getCommentReq)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, response.ResponsePostComment{
			Response: common.Response{StatusCode: -1, StatusMsg: "fail to get comments"},
			Comment:  nil,
		})
	}
	var commentListResp []response.ResponseComment

	for i := 0; i < len(commentList); i++ {
		// CommentModel转化为ResponseComment
		var responseComment response.ResponseComment

		responseComment.CommentID = int64(commentList[i].ID)
		responseComment.User = commentList[i].User
		responseComment.Content = commentList[i].Content
		responseComment.Created = commentList[i].CreatedAt.Format("2006-01-02 15:04:05")

		commentListResp = append(commentListResp, responseComment)

	}
	ctx.JSON(http.StatusOK, response.ResponseGetCommentList{
		Response:    common.Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commentListResp,
	})

}
