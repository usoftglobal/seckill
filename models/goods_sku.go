package models

import (
	"fmt"
)

type GoodsSKU struct {
	Model
	GoodsID			uint			`json:"goods_id"`
	Name 			string  	 	`json:"name"`
	Stock 			uint 		 	`json:"stock"`
	Price 			uint		 	`json:"price"`
}

func (g *GoodsSKU) CacheKey(gid uint) string {
	return fmt.Sprintf("goods_skus:%d", gid)
}