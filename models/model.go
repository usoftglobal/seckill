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
var err error

// 初始化
func init() {
	conf := libs.Conf()
	
	// MySQL 连接池
	DB, err = gorm.Open("mysql", conf.MySQL.User + ":" + conf.MySQL.Password + "@(" + conf.MySQL.Host + ")/seckill?charset=utf8&parseTime=True&loc=Local")

	// if err != nil {
	// 	假如数据库连接异常应该报警，这里不要 panic 报错因为有些业务只通过缓存访问不要影响
	// }

	DB.BlockGlobalUpdate(true)

	// Redis 连接池
	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":6379",
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

	// 虚拟字段
	CreatedAtFormat string			`gorm:"-" json:"created_at_format"`
}

// 查询之后 Hook
func (m *Model) AfterFind() (err error) {
	m.CreatedAtFormat = libs.TimeFormat(m.CreatedAt)
	return
}

// 查询之后 Hook
// 仅在组合结构覆盖 AfterFind() 时调用，类似其他语言继承里的 parent::AfterFind() 
func (m *Model) ModelAfterFind() {
	m.AfterFind()
	return 
}

// 保存之后 Hook
func (m *Model) AfterSave() (err error) {
  // Nothing ...
  return
}

// 保存之后 Hook（用法参考 ModelAfterFind）
func (m *Model) ModelAfterSave() (err error) {
  m.AfterSave()
  return
}

// 删除之后 Hook
func (m *Model) AfterDelete() (err error) {
  return
}

// 删除之后 Hook（用法参考 ModelAfterFind）
func (m *Model) ModelAfterDelete() (err error) {
	m.AfterDelete()
	return
}