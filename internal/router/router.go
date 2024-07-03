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

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Set up CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Registration endpoints for different roles
	r.POST("/register/student", studentHandler.Register)
	r.POST("/register/instructor", middleware.RBACMiddleware(), instructorHandler.Register)
	r.POST("/register/admin", middleware.RBACMiddleware(), adminHandler.Register)

	//login and logout endpoints
	r.POST("/login/student", studentHandler.Login)
	r.POST("/login/instructor", instructorHandler.Login)
	r.POST("/login/admin", adminHandler.Login)

	// r.POST("/logout", func(c *gin.Context)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.GET("/students/:id", studentHandler.GetStudentByID)
	api.GET("/students", studentHandler.GetAllStudents)
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
