package router

import (
	"awesomeProject/controller"
	"awesomeProject/middlerware"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	apiRouter := router.Group("/douyin")
	{

		// BasicRouter
		apiRouter.POST("/user/register/", controller.UserRegisterMethod)                           //用户注册
		apiRouter.POST("/user/login/", controller.UserLoginHandler)                                //验证用户登录
		apiRouter.GET("/user/", middlerware.JwtMiddleWare(), controller.GetUserInfo)               //返回用户信息
		apiRouter.GET("/feed/", middlerware.JwtMiddleWare(), controller.GetVideoList)              //请求视频流
		apiRouter.GET("/publish/list/", middlerware.JwtMiddleWare(), controller.GetPublishedVideo) //请求已发布的视频
		apiRouter.POST("/publish/action/", middlerware.JwtMiddleWare(), controller.PublishVideo)   //投稿视频

		// InteractRouter
		apiRouter.POST("/favorite/action/", middlerware.JwtMiddleWare(), controller.Favourite)     // 赞操作
		apiRouter.GET("/favorite/list/", middlerware.JwtMiddleWare(), controller.GetFavouriteList) //获取用户所有点赞视频
		apiRouter.POST("/comment/action/", middlerware.JwtMiddleWare(), controller.CommentHandler)
		apiRouter.GET("/comment/list/", middlerware.JwtMiddleWare(), controller.GetCommentList)
	}

}
