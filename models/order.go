package models

import (
	"github.com/usoftglobal/seckill/libs"
)

type Order struct {
	Model
	Uid	 			int 		 	`json:"uid"`
	OrderNo 		string 		 	`json:"order_no"`
	GoodsID 		uint 		 	`json:"goods_id"`
	Name 			string  	 	`json:"name"`
	Number 			int  		 	`json:"number"`
	Price 			int		 		`json:"price"`
	PriceCount 		int		 		`json:"price_count"`

	// 虚拟字段
	PriceFormat 	 string 		 `gorm:"-" json:"price_format"`
	PriceCountFormat string 		 `gorm:"-" json:"price_count_format"`
}

func (o *Order) AfterFind() (err error) {
	o.ModelAfterFind()
	o.PriceFormat 	   = libs.PriceFormat(o.Price)
	o.PriceCountFormat = libs.PriceFormat(o.PriceCount)
	return
}