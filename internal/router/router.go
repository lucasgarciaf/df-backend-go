package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adminsHandler "github.com/lucasgarciaf/df-backend-go/handlers/admins"
	availabilityHandler "github.com/lucasgarciaf/df-backend-go/handlers/availability"
	coursesHandler "github.com/lucasgarciaf/df-backend-go/handlers/courses"
	instructorsHandler "github.com/lucasgarciaf/df-backend-go/handlers/instructors"
	lessonsHandler "github.com/lucasgarciaf/df-backend-go/handlers/lessons"
	studentsHandler "github.com/lucasgarciaf/df-backend-go/handlers/students"
	vehiclesHandler "github.com/lucasgarciaf/df-backend-go/handlers/vehicles"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/admins"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/availability"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/courses"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/instructors"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/lessons"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/students"
	"github.com/lucasgarciaf/df-backend-go/internal/domain/vehicles"
	"github.com/lucasgarciaf/df-backend-go/internal/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(r *gin.Engine, db *mongo.Database) {
	studentRepo := students.NewMongoStudentRepository(db)
	studentService := students.NewStudentService(studentRepo)
	studentHandler := studentsHandler.NewStudentHandler(studentService)

	instructorRepo := instructors.NewMongoInstructorRepository(db)
	instructorService := instructors.NewInstructorService(instructorRepo)
	instructorHandler := instructorsHandler.NewInstructorHandler(instructorService)

	adminRepo := admins.NewMongoAdminRepository(db)
	adminService := admins.NewAdminService(adminRepo)
	adminHandler := adminsHandler.NewAdminHandler(adminService)

	courseRepo := courses.NewMongoCourseRepository(db)
	courseService := courses.NewCourseService(courseRepo)
	courseHandler := coursesHandler.NewCourseHandler(courseService)

	lessonRepo := lessons.NewMongoLessonRepository(db)
	lessonService := lessons.NewLessonService(lessonRepo)
	lessonHandler := lessonsHandler.NewLessonHandler(lessonService)

	availabilityRepo := availability.NewMongoAvailabilityRepository(db)
	availabilityService := availability.NewAvailabilityService(availabilityRepo)
	availabilityHandler := availabilityHandler.NewAvailabilityHandler(availabilityService)

	vehicleRepo := vehicles.NewMongoVehicleRepository(db)
	vehicleService := vehicles.NewVehicleService(vehicleRepo)
	vehicleHandler := vehiclesHandler.NewVehicleHandler(vehicleService)

	// Set trusted proxies
	// r.SetTrustedProxies([]string{"<your-proxy-ip-address>"})

	// Unified login and logout endpoints
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
			Role     string `json:"role" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var token string
		var err error
		switch req.Role {
		case "student":
			token, err = studentHandler.Login(c)
		// case "instructor":
		// 	token, err = instructorHandler.Login(c)
		// case "admin":
		// 	token, err = adminHandler.Login(c)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.POST("/logout", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
			Role         string `json:"role" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var err error
		switch req.Role {
		case "student":
			err = studentHandler.Logout(c)
		// case "instructor":
		// 	err = instructorHandler.Logout(c)
		// case "admin":
		// 	err = adminHandler.Logout(c)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	})

	// Registration endpoints for different roles
	r.POST("/register/student", studentHandler.Register)
	r.POST("/register/instructor", middleware.RBACMiddleware(), instructorHandler.Register)
	r.POST("/register/admin", middleware.RBACMiddleware(), adminHandler.Register)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.Use(middleware.RBACMiddleware())

	api.POST("/students", studentHandler.CreateStudent)
	api.GET("/students/:id", studentHandler.GetStudentByID)
	api.PUT("/students/:id", studentHandler.UpdateStudent)
	api.DELETE("/students/:id", studentHandler.DeleteStudent)

	api.POST("/instructors", instructorHandler.CreateInstructor)
	api.GET("/instructors/:id", instructorHandler.GetInstructorByID)
	api.PUT("/instructors/:id", instructorHandler.UpdateInstructor)
	api.DELETE("/instructors/:id", instructorHandler.DeleteInstructor)

	api.POST("/admins", adminHandler.CreateAdmin)
	api.GET("/admins/:id", adminHandler.GetAdminByID)
	api.GET("/admins/email/:email", adminHandler.GetAdminByEmail)
	api.PUT("/admins/:id", adminHandler.UpdateAdmin)
	api.DELETE("/admins/:id", adminHandler.DeleteAdmin)

	api.POST("/courses", courseHandler.CreateCourse)
	api.GET("/courses/:id", courseHandler.GetCourseByID)
	api.PUT("/courses/:id", courseHandler.UpdateCourse)
	api.DELETE("/courses/:id", courseHandler.DeleteCourse)

	api.POST("/lessons", lessonHandler.CreateLesson)
	api.GET("/lessons/:id", lessonHandler.GetLessonByID)
	api.PUT("/lessons/:id", lessonHandler.UpdateLesson)
	api.DELETE("/lessons/:id", lessonHandler.DeleteLesson)

	api.POST("/availability", availabilityHandler.CreateAvailability)
	api.GET("/availability/:id", availabilityHandler.GetAvailabilityByID)
	api.PUT("/availability/:id", availabilityHandler.UpdateAvailability)
	api.DELETE("/availability/:id", availabilityHandler.DeleteAvailability)

	api.POST("/vehicles", vehicleHandler.CreateVehicle)
	api.GET("/vehicles/:id", vehicleHandler.GetVehicleByID)
	api.PUT("/vehicles/:id", vehicleHandler.UpdateVehicle)
	api.DELETE("/vehicles/:id", vehicleHandler.DeleteVehicle)
}
