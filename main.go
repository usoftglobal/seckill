package main

import (
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/services"
)

// 应用入口
func main() {

	// 订单处理队列
	go (&services.SeckillService{}).OrderHandel()

	// GIN Framework
	ginFramework()
}

func ginFramework() {
	// gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.Default()
	ginEngine = setupRouter(ginEngine)
	ginEngine.Run(":3000")
	defer models.DB.Close()
}
