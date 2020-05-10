package services

import (
	"errors"
	"github.com/usoftglobal/seckill/models"
)

// 商品 SKU 服务
type GoodsSKUService struct {}

// 查询商品 SKU 详情
func (g *GoodsSKUService) Find(id uint) (models.GoodsSKU, error) {

	sku := models.GoodsSKU{}

	if result := models.DB.Where("goods_id = ?", id).First(&sku); result.Error != nil {
		return sku, result.Error
	}

	if sku.ID == 0 {
		return sku, errors.New("SKU 不存在或已删除")
	}

	return sku, nil
}