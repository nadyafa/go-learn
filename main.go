package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/config/db"
	"github.com/nadyafa/go-learn/controller"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/repository"
	"github.com/nadyafa/go-learn/service"
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
	userRepo := repository.NewUserRepo(dbInit)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	authRepo := repository.NewAuthRepo(dbInit)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	courseRepo := repository.NewCourseRepo(dbInit)
	couserService := service.NewCourseService(courseRepo)
	courseController := controller.NewCourseController(couserService)

	enrollRepo := repository.NewEnrollRepo(dbInit)
	enrollService := service.NewEnrollService(courseRepo, enrollRepo, userRepo)
	enrollController := controller.NewEnrollController(enrollService)

	classController := controller.NewClassController(dbInit, courseRepo)
	attendanceController := controller.NewAttendController(dbInit, courseRepo)
	projectController := controller.NewProjectController(dbInit, courseRepo)
	projectSubController := controller.NewProjectSubController(dbInit)

	// auth
	r.POST("/signup", authController.UserSignup)
	r.POST("/signin", authController.UserSignin)
	r.POST("/signout", authController.UserSignout)

	// user
	userController.GenerateAdmin()
	// admin only
	r.GET("/users", middleware.AuthMiddleware, userController.GetUsers)
	r.PUT("/users/:user_id", middleware.AuthMiddleware, userController.UpdateUserRoleByID)
	r.DELETE("/users/:user_id", middleware.AuthMiddleware, userController.DeleteUserByID)
	// admin & mentor
	r.GET("/users/:user_id", middleware.AuthMiddleware, userController.GetUserByID)

	// course
	r.POST("/courses", middleware.AuthMiddleware, courseController.CreateCourse)
	r.GET("/courses", middleware.AuthMiddleware, courseController.GetCourses)
	r.GET("/courses/:course_id", middleware.AuthMiddleware, courseController.GetCourseByID)
	r.PUT("/courses/:course_id", middleware.AuthMiddleware, courseController.UpdateCourseByID)    //admin & mentor
	r.DELETE("/courses/:course_id", middleware.AuthMiddleware, courseController.DeleteCourseByID) //admin only

	// class
	r.POST("/:course_id/classes", middleware.AuthMiddleware, classController.CreateClass) //admin & mentor
	r.GET("/:course_id/classes", middleware.AuthMiddleware, classController.GetClasses)
	r.GET("/:course_id/classes/:class_id", middleware.AuthMiddleware, classController.GetClassByID)
	r.PUT("/:course_id/classes/:class_id", middleware.AuthMiddleware, classController.UpdateClassByID)    //admin & mentor
	r.DELETE("/:course_id/classes/:class_id", middleware.AuthMiddleware, classController.DeleteClassByID) //admin & mentor

	// project
	r.POST("/:course_id/projects", middleware.AuthMiddleware, projectController.CreateProject) //admin & mentor
	r.GET("/:course_id/projects", middleware.AuthMiddleware, projectController.GetProjects)
	r.GET("/:course_id/projects/:project_id", middleware.AuthMiddleware, projectController.GetProjectByID)
	r.PUT("/:course_id/projects/:project_id", middleware.AuthMiddleware, projectController.UpdateProject)        //admin & mentor
	r.DELETE("/:course_id/projects/:project_id", middleware.AuthMiddleware, projectController.DeleteProjectByID) //admin & mentor

	// projectSub
	r.POST("/:course_id/projects/:project_id/submission", middleware.AuthMiddleware, projectSubController.StudentSubmitProject)             //student only
	r.PUT("/:course_id/projects/:project_id/submission/:project_sub_id", middleware.AuthMiddleware, projectSubController.MentorSubmitScore) //admin & mentor
	// r.GET("/:course_id/projects/:project_id/submission", middleware.AuthMiddleware, projectSubController.GetProjectSubmissions) //for all
	// r.GET("/:course_id/projects/:project_id/submission/:project_sub_id", middleware.AuthMiddleware, projectSubController.GetProjectSubmissionByID) //for all
	// r.DELETE("/:course_id/projects/:project_id/submission/:project_sub_id", middleware.AuthMiddleware, projectSubController.DeleteProjectSubmissionByID) //admin only

	// attendance
	r.POST("/:course_id/classes/:class_id/attendances", middleware.AuthMiddleware, attendanceController.StudentAttendClass)                    //admin & student
	r.GET("/:course_id/classes/:class_id/attendances", middleware.AuthMiddleware, attendanceController.GetClassAttendances)                    //admin & mentor
	r.DELETE("/:course_id/classes/:class_id/attendances/:attendance_id", middleware.AuthMiddleware, attendanceController.DeleteAttendanceByID) //admin

	// enrollment
	r.POST("/:course_id/enrollments", middleware.AuthMiddleware, enrollController.StudentEnroll)                 //student & mentor
	r.PUT("/:course_id/enrollments/:enroll_id", middleware.AuthMiddleware, enrollController.UpdateStudentEnroll) //admin

	// activity
	// percentage student attendance based on number of class
	// score of project

	// notification
	// notify mentor signup
	// notify enrollment
	// notify submission

	r.Run()
}
