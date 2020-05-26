package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usoftglobal/seckill/libs"
	"github.com/usoftglobal/seckill/services"
)

// 商品控制器
type GoodsController struct {
	GoodsService *services.GoodsService
}

func (g *GoodsController) All(c *gin.Context) {
	res, err := g.GoodsService.All()

	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}

func (g *GoodsController) Detail(c *gin.Context) {
	res, err := g.GoodsService.Find(libs.StringToUint(c.Param("id")))

	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}

func (g *GoodsController) Create(c *gin.Context) {
	err := g.GoodsService.Create(libs.StringToUint(c.Query("number")))

	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(""))
}

func (g *GoodsController) Clear(c *gin.Context) {
	res, err := g.GoodsService.Clear()

	if err != nil {
		c.JSON(http.StatusOK, libs.Fail(err))
		return
	}

	c.JSON(http.StatusOK, libs.Success(res))
}
