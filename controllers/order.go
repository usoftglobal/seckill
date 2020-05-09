package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/services"
	"github.com/usoftglobal/seckill/libs"
)

// 订单控制器
type OrderController struct {
	OrderService *services.OrderService
}

func (o *OrderController) All(c *gin.Context) {
	res, err := o.OrderService.All()
	
	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}

func (o *OrderController) Detail(c *gin.Context) {
	res, err := o.OrderService.Find(c.Param("id"))
	
	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}
