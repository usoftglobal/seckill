package main

import (
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/services"
)

// 应用入口
func main() {
	SS := new(services.SeckillService)

	go SS.OrderHandel()

	runGin()
}

func runGin() {

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.Default()
	ginEngine = setupRouter(ginEngine)
	ginEngine.Run(":3000")
	
	defer Close()
}

func Close() {
	models.DB.Close()
}
