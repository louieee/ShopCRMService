package routers

import (
	"ShopService/controllers"
	"ShopService/core"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(router *gin.Engine, db *gorm.DB) {
	// Define a route for getting a user
	baseRoute := router.Group("/api/v1/")
	allControllers := controllers.NewController(db)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	shopRoutes(baseRoute.Group("shops"), allControllers)
	userRoutes(baseRoute.Group("users"), allControllers)
	authRoutes(baseRoute.Group("auth"), allControllers)
}

func userRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	router.POST("/", controllers.RegisterUser)
	securedRoutes := router.Use(core.Auth())
	securedRoutes.GET("/", controllers.GetUser)

}
func shopRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	router.GET("/:id", controllers.GetShop)

}

func authRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	//securedRoutes := router.Use(core.Auth())
	//securedRoutes.POST("/login", controllers.GetUser)
	router.POST("/login", controllers.LoginAPI)
	router.POST("/accessToken", controllers.GetAccessTokenAPI)

}
