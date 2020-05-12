package services

import (
	"log"
	"encoding/json"
	"errors"
	"strconv"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/libs"
	"github.com/go-redis/redis/v7"
	"github.com/gin-gonic/gin"
)

var SKU *GoodsSKUService
var Order *OrderService

func init() {
	SKU   = new(GoodsSKUService)
	Order = new(OrderService)
}

type OrderData struct {
	GoodsID uint `json:"goods_id"`
	Number uint `json:"number"`
}

type SeckillService struct {}

// 购买
func (s *SeckillService) Buy(goodsID uint, numbera uint) (string, error) {
	failStr  := "fail"
	number := int64(numbera)

	if goodsID == 0 {
		return failStr, errors.New("商品ID不能为空")
	}

	// 验证购买的商品数量
	if number == 0 {
		return failStr, errors.New("至少购买一件")
	}

	// 扣减库存
	cacheKey := (&models.GoodsSKU{}).CacheKey(goodsID)
	err := s.decrStock(goodsID, number, cacheKey)
	if err != nil {
		return failStr, err
	}

	// 抢购成功 Push 到队列里
	models.RDB.LPush("GoodsMQ", libs.MapToJSON(gin.H{"goods_id": goodsID, "number": number}))

	return "success", nil
}

// 扣减库存
func (s *SeckillService) decrStock(goodsID uint, number int64, cacheKey string) error {
	err := models.RDB.Watch(func(tx *redis.Tx) error {
		// 查询当前库存
		stock, err := tx.Get(cacheKey).Int64()
		if err != nil && err != redis.Nil {
			return err
		}

		// 如果 库存 <= 0 或者 购买的数量比库存量大 返回库存不足
		if stock <= 0 || number > stock {
			return errors.New("库存不足")
		}

		// 仅在上面监听库存不变的情况下执行
		_, ExecErr := tx.TxPipelined(func(pipe redis.Pipeliner) error {
			newVal := strconv.FormatInt(stock-number, 10)
			pipe.Set(cacheKey, newVal, 0)
			return nil
		})

		return ExecErr
	}, cacheKey)

	// 如果并发情况下执行事务失败，帮助用户重新发起
	if err == redis.TxFailedErr {
		// return s.decrStock(goodsID,  number, cacheKey)

		// 这时候可能还有库存，也可以让用户自己再次发起请求，可提示抢购人数太多请重试
		return errors.New("当前抢购人数太多，请重试")
	}

	return err
}

// 异步处理订单
func (s *SeckillService) OrderHandel() {
	for {
		result, _ := models.RDB.BRPop(0, "GoodsMQ").Result()
		
		// 解析 JSON
		data := OrderData{}
		json.Unmarshal([]byte(result[1]), &data)

		// 消费队列
		err := s.updateOrder(data.GoodsID, data.Number)
		if err != nil {
			log.Println("消费错误：", err)
		}
	}
}

// 更新订单
func (s *SeckillService) updateOrder(goodsID uint, number uint) error {
	// 注意：正常情况下这里需要使用事务创建订单、同步库存

	// 生成订单
	_, err := Order.Create(goodsID, number)
	if err != nil {
		return err
	}

	// 查询商品 SKU 详情
	sku, err := SKU.Find(goodsID)
	if err != nil {
		return err
	}

	// 同步修改商品库存，这样操作不会执行 Hook 里的逻辑
	result := models.DB.Model(&sku).UpdateColumn("stock", sku.Stock - number)
	if result.Error != nil {
		return result.Error
	}

	return nil
}