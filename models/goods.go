package models

import (
	"fmt"
)

type Goods struct {
	Model
	Name 			string  	 	`json:"name"`
}

func (g *Goods) CacheKey(id uint) string {
	return fmt.Sprintf("goods:%d", id)
}