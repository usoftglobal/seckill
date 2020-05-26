package services

import (
	"errors"

	"github.com/usoftglobal/seckill/libs"
	"github.com/usoftglobal/seckill/models"
)

// 商品服务
type GoodsService struct{}

// 查询商品详情（从缓存）
func (g *GoodsService) Find(id uint) (map[string]string, error) {
	nilMap := map[string]string{}

	cacheKey := (&models.Goods{}).CacheKey(id)
	cacheResult, err := models.RDB.HGetAll(cacheKey).Result()
	if err != nil {
		return nilMap, err
	}

	if len(cacheResult) == 0 {
		return cacheResult, errors.New("商品不存在或已删除")
	}

	return cacheResult, nil
}

// 获取所有商品
func (g *GoodsService) All() ([]models.Goods, error) {
	var goods []models.Goods

	if result := models.DB.Find(&goods); result.Error != nil {
		return goods, result.Error
	}

	if len(goods) <= 0 {
		return goods, errors.New("没有数据")
	}

	return goods, nil
}

// 清空所有
func (g *GoodsService) Clear() (string, error) {
	// 清空所有数据
	models.DB.Exec("truncate table `goods`;")
	models.DB.Exec("truncate table `goods_skus`;")
	models.DB.Exec("truncate table `orders`;")
	models.RDB.FlushAll()

	// 新增
	for i := 1; i < 100; i++ {
		g.Create(1000000)
	}

	return "success", nil
}

// 创建商品
func (g *GoodsService) Create(number uint) error {

	// 开启事务
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 判断开启事务是否出错
	if err := tx.Error; err != nil {
		return err
	}

	// Goods
	goods := models.Goods{}
	goods.Name = "BMW 3系"

	if err := tx.Save(&goods).Error; err != nil {
		tx.Rollback()
		return err
	}

	// SKU
	sku := models.GoodsSKU{}
	sku.GoodsID = goods.ID
	sku.Name = "320Li"
	sku.Stock = number
	sku.Price = libs.UnitToCents(280000)

	if err := tx.Save(&sku).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Goods Cache
	goodsKey := goods.CacheKey(goods.ID)
	goodsCacheErr := models.RDB.HMSet(goodsKey, libs.StructToMap(goods)).Err()
	if goodsCacheErr != nil {
		tx.Rollback()
		return goodsCacheErr
	}

	// SKU Cache
	skuKey := sku.CacheKey(goods.ID)
	SKUCacheErr := models.RDB.Set(skuKey, sku.Stock, 0).Err()
	if SKUCacheErr != nil {
		tx.Rollback()
		return SKUCacheErr
	}

	return tx.Commit().Error
}
