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
	companyRoutes(baseRoute.Group("crm/companies"), allControllers)
	contactRoutes(baseRoute.Group("crm/contacts"), allControllers)
	leadRoutes(baseRoute.Group("crm/leads"), allControllers)
	userRoutes(baseRoute.Group("users"), allControllers)
	authRoutes(baseRoute.Group("auth"), allControllers)
}

func userRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	corsRoutes := router.Use(core.CORSMiddleware())
	securedRoutes := corsRoutes.Use(core.Auth())
	securedRoutes.GET("/me", controllers.GetUser)
	securedRoutes.GET("/", controllers.GetUserList)

}

func contactRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	securedRoutes := router.Use(core.Auth(), core.CORSMiddleware())
	//securedRoutes := router
	securedRoutes.PUT("/:contact_id", controllers.UpdateContact)
	securedRoutes.GET("/", controllers.ContactList)
	securedRoutes.POST("/", controllers.CreateContact)
	securedRoutes.DELETE("/:contact_id", controllers.DeleteContact)
	securedRoutes.GET("/:contact_id", controllers.RetrieveContact)
}

func leadRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	securedRoutes := router.Use(core.Auth(), core.CORSMiddleware())
	//securedRoutes := router
	securedRoutes.PUT("/:lead_id", controllers.UpdateLead)
	securedRoutes.GET("/", controllers.LeadList)
	securedRoutes.POST("/", controllers.CreateLead)
	securedRoutes.DELETE("/:lead_id", controllers.DeleteLead)
	securedRoutes.GET("/:lead_id", controllers.RetrieveLead)
}

func companyRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	securedRoutes := router.Use(core.Auth(), core.CORSMiddleware())
	//securedRoutes := router
	securedRoutes.PUT("/:company_id", controllers.UpdateCompany)
	securedRoutes.GET("/", controllers.CompanyList)
	securedRoutes.POST("/", controllers.CreateCompany)
	securedRoutes.DELETE("/:company_id", controllers.DeleteCompany)
	securedRoutes.GET("/:company_id", controllers.RetrieveCompany)
}

func authRoutes(router *gin.RouterGroup, controllers *controllers.DBController) {
	router.POST("/accessToken", controllers.GetAccessTokenAPI)

}
