package models

type Order struct {
	Model
	Uid	 			int 		 	`json:"uid"`
	OrderNo 		string 		 	`json:"order_no"`
	GoodsID 		uint 		 	`json:"goods_id"`
	Name 			string  	 	`json:"name"`
	Number 			uint  		 	`json:"number"`
	Price 			uint		 	`json:"price"`
	PriceCount 		uint		 	`json:"price_count"`
}