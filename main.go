package main

import (
	"payment-app/config"
	"payment-app/models"
	"payment-app/routes"
)

func main(){
	config.ConnectDatabase()
	
	config.DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Account{},
		&models.Transaction{},
	)

	r := routes.SetupRouter()

	r.Run(":8080")
	
}	