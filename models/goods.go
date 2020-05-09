package models

import (
	"fmt"
	"github.com/usoftglobal/seckill/libs"
)

type Goods struct {
	Model
	Name 			string  	 	`json:"name"`
	Stock 			int 		 	`json:"stock"`
	Price 			int		 		`json:"price"`

	// 虚拟字段
	PriceFormat 	string 		 	`gorm:"-" json:"price_format"`
}

func (g *Goods) AfterFind() (err error) {
	// 注意 ！！！
	// 修改了这部分代码，Redis 缓存里的相关字段的数据不会改变
	g.ModelAfterFind()
	g.PriceFormat = libs.PriceFormat(g.Price)
	return
}

func (g *Goods) AfterSave() (err error) {
	g.ModelAfterSave()
	g.AfterFind()

	// 插入缓存数据库
	e := RDB.HMSet(fmt.Sprintf("goods:%d", g.ID), libs.StructToMap(g)).Err()
	if e != nil {
		return e
	}

	return nil
}

func (g *Goods) AfterDelete() (err error) {
	g.ModelAfterDelete()

	// 从缓存数据库删除
	e := RDB.Del(fmt.Sprintf("goods:%d", g.ID)).Err()
	if e != nil {
		return e
	}

	return nil
}