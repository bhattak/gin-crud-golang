package main

import (
	"fmt"
	"log"
	"net/http"
	"project/model"

	"github.com/gin-gonic/gin"

	"project/database"
	"project/middleware"
	"project/service"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	// err := db.
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Initialize the middleware
	// router.Use(middleware.AuthMiddleware())

	// Initialize the user service
	userRepo := model.NewUserRepository(db)
	userAuthRepo := model.NewUserAuthRepository(db)

	userService := service.NewUserService(userRepo)
	userAuthService := service.NewUserAuthService(userAuthRepo)

	basicRoutes := router.Group("basic")
	{
		// Home route
		basicRoutes.GET("/home", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to GO Crud",
			})
		})

		basicRoutes.POST("/login", userAuthService.Login)
		basicRoutes.POST("/register", userAuthService.Register)

	}

	// Authenticated routes group
	authRoutes := router.Group("authenticated")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		// Users routes
		authRoutes.GET("/auth", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "This is authenticated route",
			})
		})
		authRoutes.GET("/fetchAll", userService.FetchAllUsers)
		authRoutes.GET("/fetchById/:id", userService.FetchUserByID)
		authRoutes.POST("/users", userService.AddUser)
		authRoutes.PUT("/users/:id", userService.UpdateUser)
		authRoutes.DELETE("/users/:id", userService.DeleteUser)
	}

	// Spin Up the server
	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

/*
	curl -X POST localhost:8080/basic/login -H 'Content-Type:application/json' -d '{"name":"ram","email":"ram@amazon.com","password":"1234"}'
	curl -X GET localhost:8080/authenticated/auth -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoicmFtQGFtYXpvbi5jb20iLCJleHAiOjE2OTEyOTI1MjV9.NFZNeY-rJIAoODFJ5ITos4qDW9JVw9fiY4MqyR0rnu4"
	curl -X DELETE localhost:8080/authenticated/users/7 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZ29wYWxAYW1hem9uLmNvbSIsImV4cCI6MTY4NjQ4MjQ5NH0.iiwX0sAzTkcKAs0FjX1oThYEew2tMs0XpHGDr2UE8tg"
*/
