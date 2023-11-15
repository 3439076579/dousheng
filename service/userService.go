package service

import (
	"awesomeProject/bloom"
	"awesomeProject/dao"
	"awesomeProject/model/user"
	"awesomeProject/model/user/request"
	"awesomeProject/redis_"
	"awesomeProject/utils"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	trumail "github.com/sdwolfe32/trumail/verifier"
	"log"
	"strconv"
	"time"
)

/*
	登录逻辑:
		1.将RequestUserLogin转化为UserLogin结构
		2.查询Redis，如果查不到执行第3步，否则验证密码正确性后返回
		3.查询数据库，验证密码正确性后返回
		4.如果登录成功，则缓存账号密码到Redis

	缓存:key:UserLogin:UserName
		val:UserID:Password
*/

// CheckUserLogin 用于验证用户登录时的账号和密码
func CheckUserLogin(userReq request.RequestUserLogin) (user.UserLogin, error) {

	var userModel user.UserLogin

	bf := bloom.NewBloomFilter(10000000, 0.03, redisService.GlobalRedis, "user_bf")

	userModel.Username = userReq.UserName
	userModel.Password = userReq.PassWord
	// 查询BloomFilter
	exist, err := bf.IsExist(userModel.Username)
	if err != nil {
		return user.UserLogin{}, err
	}
	if !exist {
		return user.UserLogin{}, utils.PassWordOrUserNameErr
	}
	// 查询Redis
	r := redisService.RedisService{
		Ctx:    context.Background(),
		PreKey: "UserLogin:",
	}
	result := r.Get(r.PreKey + userModel.Username)
	if result.Err() != redis.Nil {
		if result.Err() != nil {
			log.Println(result.Err())
			return user.UserLogin{}, result.Err()
		}
		// 验证密码正确性
		NewUserModel := utils.ParseStringToUserLogin(result.Val())
		if NewUserModel.Password == userModel.Password {
			userModel.ID = NewUserModel.ID
			return userModel, nil
		}
		return userModel, utils.PasswordIllegal
	}

	// 查询数据库
	userDao := dao.NewUserDao(false)
	NewUserModel, err := userDao.SearchUserByUserName(userModel)
	if err != nil {
		return userModel, err
	}
	// 如果密码正确，则通过校验，否则返回 PasswordIllegal错误
	if NewUserModel.Password != userModel.Password {
		return userModel, utils.PasswordIllegal
	}
	// 通过校验，把数据缓存到Redis
	r.SetNX(r.PreKey+NewUserModel.Username,
		strconv.FormatInt(NewUserModel.ID, 10)+":"+NewUserModel.Password, time.Second*3600)

	return NewUserModel, nil
}

// UserRegisterService 用于实现用户注册时的业务
func UserRegisterService(userReq *request.RequestUserLogin) (user.UserLogin, error) {

	var userModel user.UserLogin
	// 验证邮箱是否合法
	verifier := trumail.NewVerifier("haha", "wjb983798993@gmail.com")
	v, err := verifier.Verify(userReq.UserName)
	if err != nil {
		log.Println("邮箱验证出现错误")
		return user.UserLogin{}, err
	}
	if !v.Deliverable {
		return user.UserLogin{}, utils.EmailIllegal
	}

	// 邮箱合法则进行模型转换
	userModel.Username = userReq.UserName
	userModel.Password = userReq.PassWord
	userModel.ID = utils.GenerateUuid()

	// 开启事务
	userDao := dao.NewUserDao(true)
	defer userDao.DB.Commit()
	// 根据UserName查找对应记录，期望是找不到,只有报RecordNotFound错误才是期望的结果
	_, err = userDao.SearchUserByUserName(userModel)
	if err != nil && errors.Is(err, utils.RecordNotFound) {
		//如果没有通过UserName找到，再通过密码查找
		_, err = userDao.SearchUserByPassWord(userModel)
		if err != nil && errors.Is(err, utils.RecordNotFound) {
			if err = userDao.InsertUserLogin(&userModel); err != nil {
				userDao.DB.Rollback()
				return userModel, err
			}
			if err = userDao.InsertUserInfo(user.User{UserId: userModel.ID}); err != nil {
				userDao.DB.Rollback()
				return userModel, err
			}
			// 注册成功，把信息缓存到Redis
			r := redisService.RedisService{
				Ctx:    context.Background(),
				PreKey: "UserLogin:",
			}
			r.SetNX(r.PreKey+userModel.Username,
				strconv.FormatInt(userModel.ID, 10)+userModel.Password,
				time.Second*3600)
			return userModel, nil
		} else if err == nil {
			userDao.DB.Rollback()
			return userModel, utils.PassWordOrUserNameErr
		} else {
			userDao.DB.Rollback()
			return userModel, err
		}
	} else if err != nil && !errors.Is(err, utils.RecordNotFound) {
		// 数据库本身查找出现错误
		userDao.DB.Rollback()
		return userModel, err
	} else {
		// 找到了对应记录，说明该用户已存在，也是注册时的错误
		userDao.DB.Rollback()
		return userModel, utils.PassWordOrUserNameErr
	}
}

// GetUserInfoService 获取用户信息
func GetUserInfoService(userReq *request.RequestGetUserInfo) (user.User, error) {

	var userInfoModel = user.User{
		UserId: userReq.UserId,
	}

	//开启事务
	userDao := dao.NewUserDao(true)

	//userDao.SearchUserById()
	err := userDao.SearchUserInfoByUserId(&userInfoModel)
	if err != nil {
		userDao.DB.Rollback()
		return userInfoModel, err
	}
	userDao.DB.Commit()
	return userInfoModel, nil
}
