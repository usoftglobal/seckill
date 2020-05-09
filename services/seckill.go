package services

import (
	"log"
	"errors"
	"strconv"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
)

var GS *GoodsService
var OS *OrderService

func init() {
	GS = new(GoodsService)
	OS = new(OrderService)
}

type SeckillService struct {}

// 购买
func (s *SeckillService) Buy(goodsID string, number string) (string, error) {
	
	failStr := "fail"

	// 验证购买的商品数量
	if number == "" {
		return failStr, errors.New("商品数量不能为空")
	}

	int64Number, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return failStr, err
	}

	if int64Number < 1 {
		return failStr, errors.New("至少购买一个")
	}

	// 1、从缓存里查出商品详情
	detail, err := GS.FindFromCache(goodsID)
	
	if err != nil {
		return failStr, err
	}

	// 2、检查商品库存是否足够这个人购买的
	if number > detail["stock"] {
		
		// 扣除相对应的库存
		result, err := models.RDB.HIncrBy("goods:" + goodsID, "stock", -int64Number).Result()
		if err != nil {
			return failStr, err
		}

		log.Println("日志：", number, detail["stock"], result)

		// Push 到队列里
		d := map[string]interface{}{"goods_id": goodsID, "number": int64Number}
		models.RDB.LPush("GoodsMQ", libs.MapToJSON(d))

		return "success", nil
	} else {
		return failStr, errors.New("库存不足")
	}
}

// 异步处理订单
func (s *SeckillService) OrderHandel() {
	for {
		result, _ := models.RDB.BRPop(0, "GoodsMQ").Result()
		resultMap := libs.JSONToMap(result[1])
		s.updateOrder(resultMap["goods_id"].(string), int(resultMap["number"].(float64)))
	}
}

// 更新订单
func (s *SeckillService) updateOrder(goodsID string, number int) error {
	// 注意：正常情况下这里需要使用事务创建订单、同步库存

	// 生成订单
	_, err := OS.Create(goodsID, number)
	if err != nil {
		return err
	}

	// 查询商品详情
	goods, err := GS.Find(goodsID)
	if err != nil {
		return err
	}

	// 同步修改商品库存，这样操作不会执行 Hook 里的逻辑
	result := models.DB.Model(&goods).UpdateColumn("stock", goods.Stock - number)
	if result.Error != nil {
		return result.Error
	}

	return nil
}