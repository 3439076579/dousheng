package controller

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user/request"
	"awesomeProject/model/user/response"
	"awesomeProject/service"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// UserLoginHandler 用户登录
func UserLoginHandler(router *gin.Context) {

	var userLogin request.RequestUserLogin

	userLogin.UserName = router.Query("username")
	userLogin.PassWord = router.Query("password")

	userLoginModel, err := service.CheckUserLogin(userLogin)
	if err != nil {
		log.Println("验证出现错误", err)
		router.JSON(http.StatusOK, &response.ResponseUserLogin{
			Response: common.Response{StatusCode: -1, StatusMsg: "验证出现错误"},
		})
		return
	}

	j := utils.NewJwt()
	token, err := j.CreateToken(userLoginModel.Username, userLoginModel.ID)
	if err != nil {
		log.Println("获取token失败", err)
		//如果校验失败，则返回请求失败的报文
		router.JSON(http.StatusOK, &response.ResponseUserLogin{
			Response: common.Response{StatusCode: -1, StatusMsg: "服务器出现错误，请重试"},
		})
		return
	}

	router.JSON(http.StatusOK, &response.ResponseUserLogin{
		Response: common.Response{StatusCode: 0, StatusMsg: "success"},
		Token:    token,
		UserID:   userLoginModel.ID,
	})
}

// UserRegisterMethod 用户注册
func UserRegisterMethod(router *gin.Context) {

	var userLogin request.RequestUserLogin
	//获取用户信息
	userLogin.UserName = router.Query("username")
	userLogin.PassWord = router.Query("password")
	/*根据用户信息进行校验，成功则返回成功报文（调用service层方法实现）
	有两种抛出异常的情况
	一：邮箱地址不合法，二、账号或密码已存在（合并账号存在或密码存在两种情况）*/
	userModel, err := service.UserRegisterService(&userLogin)
	if err != nil {
		fmt.Println("出现错误", err)
		router.JSON(http.StatusOK, &response.ResponseUserLogin{
			Response: common.Response{StatusCode: -1, StatusMsg: "出现错误"},
			Token:    "",
			UserID:   -1,
		})
		return
	}

	var token string
	j := utils.NewJwt()
	//根据用户名产生token，
	token, err = j.CreateToken(userModel.Username, userModel.ID)
	if err != nil {
		fmt.Println("获取token失败")
		router.JSON(http.StatusOK, &response.ResponseUserLogin{
			Response: common.Response{StatusCode: 0, StatusMsg: "获取token失败"},
			Token:    "",
			UserID:   -1,
		})
		return
	}

	router.JSON(http.StatusOK, &response.ResponseUserLogin{
		Response: common.Response{StatusCode: 0},
		Token:    token,
		UserID:   userModel.ID,
	})

}

// GetUserInfo 返回用户信息
func GetUserInfo(router *gin.Context) {

	var user request.RequestGetUserInfo
	id, err := strconv.Atoi(router.Query("user_id"))
	user.UserId = int64(id)
	userInfoModel, err := service.GetUserInfoService(&user)
	if err != nil {
		fmt.Println("获取用户信息出错")
		router.JSON(http.StatusOK, &response.ResponseGetUserInfo{
			Response: common.Response{
				StatusCode: -1,
				StatusMsg:  "获取用户信息出错",
			},
		})
		return
	}

	router.JSON(http.StatusOK, &response.ResponseGetUserInfo{
		Response: common.Response{StatusCode: 0},
		User:     userInfoModel,
	})
}
