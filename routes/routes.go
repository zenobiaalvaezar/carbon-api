package routes

import (
	"carbon-api/caches"
	"carbon-api/config"
	"carbon-api/controllers"
	"carbon-api/middlewares"
	"carbon-api/repositories"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	// Fuel routes
	fuelRepository := repositories.NewFuelRepository(config.DB)
	fuelCache := caches.NewFuelCache(config.RedisClient)
	fuelController := controllers.NewFuelController(fuelRepository, fuelCache)

	f := e.Group("/fuels")
	// TODO: add check auth & check user admin
	f.GET("", fuelController.GetAllFuels)
	f.GET("/:id", fuelController.GetFuelByID)
	f.POST("", fuelController.CreateFuel)
	f.PUT("/:id", fuelController.UpdateFuel)
	f.DELETE("/:id", fuelController.DeleteFuel)

	// Carbon fuel routes
	carbonFuelRepository := repositories.NewCarbonFuelRepository(config.DB)
	carbonFuelController := controllers.NewCarbonFuelController(carbonFuelRepository)

	cf := e.Group("/carbon-fuels")
	// TODO: add check auth & check user customer
	cf.GET("", carbonFuelController.GetAllCarbonFuels)
	cf.GET("/:id", carbonFuelController.GetCarbonFuelByID)
	cf.POST("", carbonFuelController.CreateCarbonFuel)
	cf.DELETE("/:id", carbonFuelController.DeleteCarbonFuel)

	// Carbon summary routes
	carbonSummaryRepository := repositories.NewCarbonSummaryRepository(config.DB)
	carbonSummaryController := controllers.NewCarbonSummaryController(carbonSummaryRepository)

	cs := e.Group("/carbon-summaries")
	// TODO: add check auth & check user customer
	cs.GET("", carbonSummaryController.GetCarbonSummary)

	// Role routes
	roleRepository := repositories.NewRoleRepository(config.DB)
	roleController := controllers.NewRoleController(roleRepository)

	r := e.Group("/roles")
	r.GET("", roleController.GetAllRoles)
	r.GET("/:id", roleController.GetRoleByID)
	r.POST("", roleController.CreateRole)
	r.PUT("/:id", roleController.UpdateRole)
	r.DELETE("/:id", roleController.DeleteRole)

	// Initialize repositories and controllers
	userRepository := repositories.NewUserRepository(config.DB)
	userController := controllers.NewUserController(userRepository)

	// Public routes for user
	e.POST("/register", userController.RegisterUser)
	e.POST("/login", userController.LoginUser)

	// Protected routes for user profile
	userGroup := e.Group("/users")
	userGroup.Use(middlewares.CheckAuth) // Gunakan middleware CheckAuth
	userGroup.GET("/profile", userController.GetProfile)
	userGroup.Use(middlewares.CheckAuth)
	userGroup.POST("/logout", userController.LogoutUser)

	userGroup.Use(middlewares.CheckAuth)
	userGroup.PUT("/profile", userController.UpdateProfile)

	// electric
	electricRepository := repositories.NewElectricRepository(config.DB)
	electricCache := caches.NewElectricCache(config.RedisClient)
	electricController := controllers.NewElectricController(electricRepository, electricCache)

	l := e.Group("/electrics")
	l.GET("", electricController.GetAllElectrics)
	l.GET("/:id", electricController.GetElectricByID)
	l.POST("", electricController.CreateElectric)
	l.PUT("/:id", electricController.UpdateElectric)
	l.DELETE("/:id", electricController.DeleteElectric)
}
