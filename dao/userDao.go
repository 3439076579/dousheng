package dao

import (
	"awesomeProject/model/user"
	"awesomeProject/utils"
	"gorm.io/gorm"
)

type UserDao struct {
	DB *gorm.DB
}

func NewUserDao(UserTransaction bool) UserDao {
	return UserDao{DB: utils.GetDB(UserTransaction)}
}
func NewUserDaoFromDB(db *gorm.DB) UserDao {
	return UserDao{DB: db}
}

// SearchUserById 查询UserLogin通过ID
func (u UserDao) SearchUserById(login user.UserLogin) (user.UserLogin, error) {

	var userModel user.UserLogin

	res := u.DB.Model(&user.UserLogin{}).Where(
		"user_id=?", login.ID).Find(&userModel)
	if res.Error != nil {
		return userModel, res.Error
	} else {
		if res.RowsAffected == 0 {
			return userModel, utils.RecordNotFound
		}
	}

	return userModel, nil
}
func (u UserDao) SearchUserByUserName(login user.UserLogin) (user.UserLogin, error) {

	var userModel user.UserLogin

	res := u.DB.Model(&user.UserLogin{}).Where(
		"username=?", login.Username).Find(&userModel)
	if res.Error != nil {
		return userModel, res.Error
	} else {
		if res.RowsAffected == 0 {
			return userModel, utils.RecordNotFound
		}
	}
	return userModel, nil
}
func (u UserDao) SearchUserByPassWord(login user.UserLogin) (user.UserLogin, error) {

	var userModel user.UserLogin
	res := u.DB.Model(&user.UserLogin{}).Where(
		"password=?", login.Password).Find(&userModel)
	if res.Error != nil {
		return userModel, res.Error
	} else {
		if res.RowsAffected == 0 {
			return userModel, utils.RecordNotFound
		}
	}
	return userModel, nil
}
func (u UserDao) InsertUserLogin(login *user.UserLogin) error {
	res := u.DB.Model(&user.UserLogin{}).Create(
		login,
	)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (u UserDao) SearchUserInfoById(userModel *user.User) error {

	res := u.DB.Model(&user.User{}).Where(
		"id=?", userModel.ID,
	).Find(&userModel)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil
}
func (u UserDao) SearchUserInfoByUserId(userModel *user.User) error {

	res := u.DB.Model(&user.User{}).Where(
		"user_id=?", userModel.UserId,
	).Find(&userModel)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}
	return nil
}

func (u UserDao) InsertUserInfo(userArr ...user.User) error {
	// 如果传递进来的可变参数数组长度不为0，那么就不使用默认值
	if len(userArr) > 1 {
		return utils.ParameterErr
	}
	// 即userArr有值，不使用默认值
	if len(userArr) > 0 {
		user_ := &user.User{
			UserId: userArr[0].UserId,
		}
		if res := u.DB.
			Model(&user.User{}).
			Create(user_); res.Error != nil {
			return res.Error
		}
	} else { // userArr长度为0，使用默认值
		if res := u.DB.
			Model(&user.User{}).
			Create(&user.User{UserId: utils.GenerateUuid()}); res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (u UserDao) IncreaseWorkCount(UserModel user.User, IncrNum int64) error {

	res := u.DB.Model(&user.User{}).
		Where("id=?", UserModel.ID).
		Update("work_count", gorm.Expr("work_count+?", IncrNum))
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}

	return nil
}

func (u UserDao) IncreaseFavoriteCount(UserModel user.User) error {

	res := u.DB.Model(&user.User{}).
		Where("user_id=?", UserModel.UserId).
		Update("favorite_count", UserModel.FavoriteCount)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return utils.RecordNotFound
		}
	}

	return nil
}
