package services

import (
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
)

// 订单服务
type OrderService struct {}

// 创建订单
func (o *OrderService) Create(goodsID uint, number uint) (string, error) {

	orderNo := libs.CreateOrderNo()
	order 	:= models.Order{}
	
	// 查询商品详情
	GoodsService := new(GoodsService)
	goods, err := GoodsService.FindFromDB(goodsID)
	if err != nil {
		return orderNo, err
	}

	// 查询 SKU 详情
	SKUService := new(GoodsSKUService)
	sku, err := SKUService.Find(goodsID)
	if err != nil {
		return orderNo, err
	}

	// 创建订单
	order.Uid 	 	  = libs.CreateRandNo()
	order.OrderNo 	  = orderNo
	order.Name    	  = goods.Name
	order.GoodsID 	  = goods.ID
	order.Number  	  = number
	order.Price   	  = sku.Price
	order.PriceCount  = sku.Price * number

	if result := models.DB.Save(&order); result.Error != nil {
		return orderNo, result.Error
	}

	return orderNo, nil
}