package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/usoftglobal/seckill/libs"
	"github.com/usoftglobal/seckill/models"
)

// 订单服务
type OrderService struct{}

// 创建订单处理
func (o *OrderService) CreateOrderHandel(goodsID uint, number uint) error {
	// 事务如果返回不为 nil 整个事务都会 rollback
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 查询\判断商品详情
		goods, err := o.findAndCheckGoods(tx, goodsID)
		if err != nil {
			return err
		}

		// 查询\判断库存
		sku, err := o.findAndCheckStock(tx, goodsID)
		if err != nil {
			return err
		}

		// 扣减库存
		if err := tx.Model(&sku).UpdateColumn("stock", gorm.Expr("stock - ?", number)).Error; err != nil {
			return err
		}

		// 创建订单
		_, err = o.createOrder(tx, goods, sku, number)
		if err != nil {
			return err
		}

		return nil // 返回 nil 提交事务
	})
}

// 查询\判断商品详情
func (o *OrderService) findAndCheckGoods(tx *gorm.DB, id uint) (models.Goods, error) {

	goods := models.Goods{}

	if id == 0 {
		return goods, errors.New("商品ID不能为空")
	}

	if result := tx.First(&goods, id); result.Error != nil {
		return goods, result.Error
	}

	if goods.ID == 0 {
		return goods, errors.New("商品不存在或已删除")
	}

	return goods, nil
}

// 查询\判断库存
func (o *OrderService) findAndCheckStock(tx *gorm.DB, id uint) (models.GoodsSKU, error) {

	sku := models.GoodsSKU{}

	if result := tx.Where("goods_id = ?", id).First(&sku); result.Error != nil {
		return sku, result.Error
	}

	if sku.ID == 0 {
		return sku, errors.New("SKU 不存在或已删除")
	}

	if sku.Stock <= 0 {
		return sku, errors.New("库存不足")
	}

	return sku, nil
}

// 创建订单
func (o *OrderService) createOrder(tx *gorm.DB, goods models.Goods, sku models.GoodsSKU, number uint) (string, error) {
	order := models.Order{}
	orderNo := libs.CreateOrderNo()

	// 创建订单
	order.Uid = libs.CreateRandNo()
	order.OrderNo = orderNo
	order.Name = goods.Name
	order.GoodsID = goods.ID
	order.Number = number
	order.Price = sku.Price
	order.PriceCount = sku.Price * number

	if result := tx.Save(&order); result.Error != nil {
		return orderNo, result.Error
	}

	return orderNo, nil
}
