package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/config/db"
	"github.com/nadyafa/go-learn/controller"
	"github.com/nadyafa/go-learn/middleware"
)

func main() {
	// initiate db
	dbInit, err := db.DBInit()
	if err != nil {
		log.Fatalf("Unable initializing DB: %v", err)
	}

	// check db connection
	dbConn, err := dbInit.DB()
	if err != nil {
		log.Fatalf("Unable connect to DB: %v", err)
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Unable ping DB connection: %v", err)
	}

	// db migration
	db.RunMigration(dbInit)

	// create super admin
	// if err := model.AdminLogin(dbInit); err != nil {
	// 	log.Printf("Unable creating super admin: %v", err)
	// }

	// setup route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// setup dependencies injection
	// userRepo := repository.NewUserRepo(dbInit)
	// userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(dbInit)
	adminController := controller.NewAdminController(dbInit)
	courseController := controller.NewCourseController(dbInit)

	// auth
	r.POST("/signup", userController.UserSignup)
	r.POST("/signin", userController.UserSignin)
	r.POST("/signout", userController.UserSignout)

	// user
	adminController.GenerateAdmin()
	// admin only
	r.GET("/users", middleware.AuthMiddleware, adminController.GetUsers)
	r.PUT("/users/:user_id", middleware.AuthMiddleware, adminController.UpdateUserRoleByID)
	r.DELETE("/users/:user_id", middleware.AuthMiddleware, adminController.DeleteUserByID)
	// admin & mentor
	r.GET("/users/:user_id", middleware.AuthMiddleware, adminController.GetUserByID)

	// course
	r.POST("/courses", middleware.AuthMiddleware, courseController.CreateCourse)
	// r.GET("/courses", courseController.GetAllCourses)
	// r.GET("/courses/:course_id", courseController.GetCourseByID)
	// r.PUT("/courses/:course_id", controller.AuthMiddleware, courseController.UpdateCourse) //admin & mentor
	// r.DELETE("/courses/:user_id", controller.AuthMiddleware, courseController.DeleteCourseByID) //admin only

	r.Run()
}
