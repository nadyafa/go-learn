package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/config/db"
	"github.com/nadyafa/go-learn/config/helper"
	"github.com/nadyafa/go-learn/controller"
	"github.com/nadyafa/go-learn/repository"
	"github.com/nadyafa/go-learn/service"
)

func main() {
	// initiate db
	dbInit, err := db.DBInit()
	if err != nil {
		log.Fatal("Error initializing DB:", err)
	}

	// check db connection
	dbConn, err := dbInit.DB()
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, "DB connection error", err)
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		helper.Logger(helper.LoggerLevelError, "Unable to ping DBConn", err)
	}

	// db migration
	db.RunMigration(dbInit)

	// setup route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// setup dependencies
	userRepo := repository.NewUserRepo(dbInit)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// users
	r.POST("/users", userController.UserSignup)

	r.Run()
}
