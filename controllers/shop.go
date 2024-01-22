package controllers

import (
	_ "ShopService/core"
	"ShopService/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (dc *DBController) GetShop(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	if err := dc.DB.Preload("Shops").First(&user, userID).Error; err != nil {

		return
	}
	c.JSON(http.StatusOK, user)
}
