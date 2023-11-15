package middlerware

import (
	"awesomeProject/model/common"
	"awesomeProject/model/user"
	"awesomeProject/model/user/response"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var mutex sync.Mutex

// JwtMiddleWare 因为token的位置不一样，而且对于不同token来说，对token的需要程度不同
// 有些接口依赖于token里面解析出来的user_id，例如，获取发布列表和发布视频的接口
func JwtMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		url := c.Request.URL
		var token string

		switch url.Path {
		// 对于请求参数中含有token的执行如下的分支
		case "/douyin/user/", "/douyin/publish/list/", "/douyin/favorite/action/",
			"/douyin/favorite/list/", "/douyin/feed/", "/douyin/comment/action/",
			"/douyin/comment/list/":
			token = c.Query("token")
		case "douyin/publish/action/":
			token = c.PostForm("token")
		}
		// 执行完上述分支，token还是空值，那么只有可能是不存在token，拒绝请求
		if token == "" {
			if url.Path == "/douyin/feed/" {
				c.Next()
			}
			c.Abort()
			c.JSON(http.StatusOK, &common.Response{
				StatusCode: -1,
			})
			return
		}
		j := utils.NewJwt()

		//解析token信息，处理鉴权问题，解析token时出现错误就
		claims, err := j.ParseToken(token)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, &response.ResponseGetUserInfo{
				Response: common.Response{
					StatusCode: -1,
					StatusMsg:  "解析token时产生错误",
				},
				User: user.User{},
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
