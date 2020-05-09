package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/services"
	"github.com/usoftglobal/seckill/libs"
)

// 秒杀控制器
type SeckillController struct {
	SeckillService *services.SeckillService
}

func (s *SeckillController) Buy(c *gin.Context) {
	res, err := s.SeckillService.Buy(c.Query("goods_id"), c.Query("number"))
	
	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}