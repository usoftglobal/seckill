package models

import (
	"time"
  	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/usoftglobal/seckill/libs"
	"github.com/go-redis/redis/v7"
)

var DB *gorm.DB
var RDB *redis.Client

// 初始化
func init() {
	conf := libs.Conf()
	
	// MySQL 连接池
	dsn := conf.MySQL.User + ":" + conf.MySQL.Password + "@(" + conf.MySQL.Host + ")/" + conf.MySQL.DB
	DB, _ = gorm.Open("mysql", dsn + "?charset=utf8&parseTime=True&loc=Local")

	// if err != nil {
	// 	假如数据库连接异常应该报警，这里不要 panic 报错因为有些业务只通过缓存访问不要影响
	// }

	DB.BlockGlobalUpdate(true)

	// Redis 连接池
	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":6399",
		Password: "",
		DB:       0,
	})
}

// 基础数据库模型
type Model struct {
	ID				uint		 	`json:"id"`
	CreatedAt 		time.Time	 	`json:"created_at"`
	UpdatedAt 		time.Time	 	`json:"updated_at"`
	DeletedAt 		*time.Time	 	`json:"deleted_at"`
}