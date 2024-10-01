package routes

import (
	"payment-app/controllers"
	"payment-app/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    user := router.Group("/user")
    {
        user.POST("/register", controllers.RegisterUser)
        user.POST("/login", controllers.LoginUser)
    }

    transaction := router.Group("/transaction")
    transaction.Use(middleware.JwtAuthMiddleware())
    {
        transaction.POST("/create", controllers.CreateTransaction)
        transaction.GET("/history/:userID", controllers.GetTransactionHistory)
    }

    account := router.Group("/account")
    account.Use(middleware.JwtAuthMiddleware())
    {
        account.POST("/add", controllers.AddAccount)
    }

    wallet := router.Group("/wallet")
    wallet.Use(middleware.JwtAuthMiddleware())
    {
        wallet.POST("/create", controllers.CreateWallet) 
        wallet.GET("/:userID", controllers.GetWallet)
        wallet.POST("/add-balance", controllers.AddBalance) 
    }

    return router
}
