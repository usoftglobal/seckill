package services

import (
	"errors"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
)

// 订单服务
type OrderService struct {}

func (o *OrderService) All() ([]models.Order, error) {
	var order []models.Order
	
	if result := models.DB.Find(&order); result.Error != nil {
		return order, result.Error
	}

	if len(order) <= 0 {
		return order, errors.New("没有数据")
	}

	return order, nil
}

func (o *OrderService) Find(id string) (models.Order, error) {

	order := models.Order{}
	
	if result := models.DB.First(&order, id); result.Error != nil {
		return order, result.Error
	}

	if order.ID == 0 {
		return order, errors.New("订单不存在或已删除")
	}

	return order, nil
}

func (o *OrderService) Create(goodsID string, number int) (string, error) {

	orderNo := libs.CreateOrderNo()
	order 	:= models.Order{}
	
	// 查询商品详情
	GoodsService := new(GoodsService)
	goods, err := GoodsService.Find(goodsID)
	if err != nil {
		return orderNo, err
	}

	// 创建订单
	order.Uid 	 	  = libs.CreateRandNo()
	order.OrderNo 	  = orderNo
	order.Name    	  = goods.Name
	order.GoodsID 	  = goods.ID
	order.Number  	  = number
	order.Price   	  = goods.Price
	order.PriceCount  = goods.Price * number

	if result := models.DB.Save(&order); result.Error != nil {
		return orderNo, result.Error
	}

	return orderNo, nil
}