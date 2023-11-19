package utils

import (
	"awesomeProject/model/interactor"
	"awesomeProject/model/user"
	"awesomeProject/model/video"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"reflect"
	"time"
)

var (
	globalDB *gorm.DB
)

type TransactionsKey struct{}

func InitDB() {
	dsn := "root:123456@tcp(localhost:3306)/douyin_projoect?charset=utf8&parseTime=True&loc=Local"
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
		})
		if err != nil {
			fmt.Println("连接数据库出现错误", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if db != nil {
			break
		}
	}
	if db == nil {
		panic("连接数据库失败")
	}

	s, err := db.DB()
	s.SetMaxOpenConns(10)
	s.SetMaxIdleConns(5)
	s.SetConnMaxLifetime(time.Hour)
	log.Println("MySQL Connect Successed")
	// 进行了四个表的自动迁移
	res := db.AutoMigrate(&user.User{}, &user.UserLogin{}, &video.Video{}, &interactor.Comment{})
	if res != nil {
		panic(res)
		return
	}
	log.Println("Table Automigrate Successed")
	globalDB = db
}

func GetDB(UseTransaction bool) *gorm.DB {

	if UseTransaction == true {
		ctx := begintransaction()
		value := ctx.Value(TransactionsKey{})
		if value != nil {
			tx, ok := value.(*gorm.DB)
			if !ok {
				log.Panic("unexpect context value type:", reflect.TypeOf(tx))
				return nil
			}
			return tx
		} else {
			log.Panic("fail to get Transaction")
			return nil
		}
	} else {
		return globalDB.WithContext(context.Background())
	}

}

func begintransaction() context.Context {
	ctx := context.Background()
	tx := globalDB.WithContext(ctx).Begin()
	ctx = context.WithValue(ctx, TransactionsKey{}, tx)
	return ctx
}
