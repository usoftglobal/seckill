package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/controllers"
)

func setupRouter(r *gin.Engine) *gin.Engine {

	goods 	:= new(controllers.GoodsController)
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
	r.GET("/goodsUpdate/:id", goods.Update)
	r.GET("/goodsDelete/:id", goods.Delete)

	// 秒杀
	r.GET("/seckill/buy", seckill.Buy)

	return r
}