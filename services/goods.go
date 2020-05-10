package services

import (
	"errors"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
)

// 商品服务
type GoodsService struct {}

// 查询商品详情
func (g *GoodsService) Find(id uint) (map[string]string, error) {
	nilMap := map[string]string {}

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
	sku.Name    = "320Li"
	sku.Stock   = number
	sku.Price   = libs.UnitToCents(280000)
	
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

// ---------------------------------- 下面的不重要 ---------------------------------- //

// 查询商品详情从数据库
func (g *GoodsService) FindFromDB(id uint) (models.Goods, error) {

	goods := models.Goods{}

	if id == 0 {
		return goods, errors.New("商品ID不能为空")
	}

	if result := models.DB.First(&goods, id); result.Error != nil {
		return goods, result.Error
	}

	if goods.ID == 0 {
		return goods, errors.New("商品不存在或已删除")
	}

	return goods, nil
}

// 修改商品
func (g *GoodsService) Update(id string) (bool, error) {
	// current, err := g.Find(id)
	
	// if err != nil {
	// 	return false, err
	// }

	// data := map[string]interface{}{"name": "NewName"}

	// if result := models.DB.Model(&current).Updates(data); result.Error != nil {
	// 	return false, result.Error
	// }

	return true, nil
}

// 删除商品
func (g *GoodsService) Delete(id string) (bool, error) {
	// current, err := g.Find(id)
	
	// if err != nil {
	// 	return false, err
	// }
 
	// if result := models.DB.Delete(&current); result.Error != nil {
	// 	return false, result.Error
	// }

	return true, nil
}