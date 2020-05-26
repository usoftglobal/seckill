package main

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/controllers"
	"github.com/usoftglobal/seckill/models"
	"github.com/usoftglobal/seckill/services"
)

// 应用入口
func main() {
	app()
}

// 秒杀应用
func app() {
	// 订单处理队列
	go (&services.SeckillService{}).OrderHandel()

	// GIN Framework
	ginFramework()
}

func ginFramework() {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	pprof.Register(ginEngine, "dd/pprof")
	ginEngine = setupRouter(ginEngine)
	ginEngine.Run(":3000")

	defer models.DB.Close()
}

func setupRouter(r *gin.Engine) *gin.Engine {

	goods := new(controllers.GoodsController)
	seckill := new(controllers.SeckillController)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Golang!")
	})

	// 清空
	r.GET("/clear", goods.Clear)

	// 商品
	r.GET("/goods", goods.All)
	r.GET("/goods/:id", goods.Detail)
	r.GET("/goodsCreate", goods.Create)

	// 秒杀
	r.GET("/seckill/buy", seckill.Buy)

	return r
}
