package services

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/libs"
	"github.com/usoftglobal/seckill/models"
)

type OrderData struct {
	GoodsID uint `json:"goods_id"`
	Number  uint `json:"number"`
}

type SeckillService struct{}

// 购买
func (s *SeckillService) Buy(goodsID uint, unitNumber uint) (string, error) {
	failStr := "fail"
	number := int64(unitNumber)

	// 购买检查
	err := s.buyCheck(goodsID, number)
	if err != err {
		return failStr, err
	}

	// 扣减库存
	cacheKey := (&models.GoodsSKU{}).CacheKey(goodsID)
	err = s.decrStockBy(number, cacheKey)
	if err != nil {
		return failStr, err
	}

	// 抢购成功 Push 到队列里
	models.RDB.LPush("GoodsMQ", libs.MapToJSON(gin.H{"goods_id": goodsID, "number": number}))

	return "success", nil
}

// 购买检查
func (s *SeckillService) buyCheck(goodsID uint, number int64) error {
	if goodsID == 0 {
		return errors.New("商品ID不能为空")
	}

	// 验证购买的商品数量
	if number == 0 {
		return errors.New("至少购买一件")
	}

	return nil
}

// 扣减库存
func (s *SeckillService) decrStockBy(number int64, cacheKey string) error {

	var nilStatus int64 = 0

	luaScript := `
		local key    = tostring(KEYS[1])
		local number = tonumber(ARGV[1])
		
		-- 库存判断
		local stock  = redis.call("GET", key)
		stock		 = tonumber(stock)
		
		if(number > stock) then
			return 0
		end
		
		-- 库存扣减
		redis.call("DECRBY", key, number)

		return 1
	`

	// 将脚本生成一个 sha1 哈希值，减少网络传输
	luaScriptSha, err := models.RDB.ScriptLoad(luaScript).Result()
	if err != nil {
		return err
	}

	// 执行脚本
	status := models.RDB.EvalSha(luaScriptSha, []string{cacheKey}, number)

	if status.Val() == nilStatus {
		return errors.New("库存不足")
	}

	return nil
}

// 异步处理订单
func (s *SeckillService) OrderHandel() {
	Order := new(OrderService)

	for {
		result, _ := models.RDB.BRPop(0, "GoodsMQ").Result()

		// 解析 JSON
		data := OrderData{}
		json.Unmarshal([]byte(result[1]), &data)

		// 消费队列
		err := Order.CreateOrderHandel(data.GoodsID, data.Number)
		if err != nil {
			log.Println("消费错误：", err)
			// 正常这里应该再回滚到队列里进行重试处理
		}
	}
}
