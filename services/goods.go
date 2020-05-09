package services

import (
	"fmt"
	"errors"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
)

// 商品服务
type GoodsService struct {}

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

// 查询商品详情
func (g *GoodsService) Find(id string) (models.Goods, error) {

	goods := models.Goods{}

	if id == "" {
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

// 查询商品详情从缓存
func (g *GoodsService) FindFromCache(id string) (map[string]string, error) {
	m := map[string]string {}

	cacheResult, err := models.RDB.HGetAll(fmt.Sprintf("goods:%s", id)).Result()
	if err != nil {
		return m, err
	}

	if len(cacheResult) == 0 {
		return cacheResult, errors.New("商品不存在或已删除")
	}

	return cacheResult, nil
}

// 创建商品
func (g *GoodsService) Create() (models.Goods, error) {

	goods := models.Goods{}

	name := "康师傅冰红茶"
	stock:= 10000
	price:= libs.UnitToCents("2.2")

	goods.Name  = name
	goods.Stock = stock
	goods.Price = price

	if result := models.DB.Save(&goods); result.Error != nil {
		return goods, result.Error
	}

	return g.Find(fmt.Sprintf("%d", goods.ID))
}

// 修改商品
func (g *GoodsService) Update(id string) (bool, error) {
	current, err := g.Find(id)
	
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{"name": "崔", "stock": 0}

	if result := models.DB.Model(&current).Updates(data); result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// 删除商品
func (g *GoodsService) Delete(id string) (bool, error) {
	current, err := g.Find(id)
	
	if err != nil {
		return false, err
	}
 
	if result := models.DB.Delete(&current); result.Error != nil {
		return false, result.Error
	}

	return true, nil
}