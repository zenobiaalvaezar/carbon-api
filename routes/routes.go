package routes

import (
	"carbon-api/caches"
	"carbon-api/config"
	"carbon-api/controllers"
	"carbon-api/middlewares"
	"carbon-api/repositories"

	_ "carbon-api/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(e *echo.Echo) {
	// Role routes
	roleRepository := repositories.NewRoleRepository(config.DB)
	roleController := controllers.NewRoleController(roleRepository)

	r := e.Group("/roles")
	r.Use(middlewares.CheckAuth)
	r.Use(middlewares.CheckRoleAdmin)
	r.GET("", roleController.GetAllRoles)
	r.GET("/:id", roleController.GetRoleByID)
	r.POST("", roleController.CreateRole)
	r.PUT("/:id", roleController.UpdateRole)
	r.DELETE("/:id", roleController.DeleteRole)

	// User routes
	userRepository := repositories.NewUserRepository(config.DB)
	userController := controllers.NewUserController(userRepository)

	e.POST("/register", userController.RegisterUser)
	e.POST("/login", userController.LoginUser)

	userGroup := e.Group("/users")
	userGroup.Use(middlewares.CheckAuth)
	userGroup.GET("/profile", userController.GetProfile)
	userGroup.POST("/logout", userController.LogoutUser)
	userGroup.PUT("/profile", userController.UpdateProfile)
	userGroup.PUT("/password", userController.UpdatePassword)

	// Fuel routes
	fuelRepository := repositories.NewFuelRepository(config.DB)
	fuelCache := caches.NewFuelCache(config.RedisClient)
	fuelController := controllers.NewFuelController(fuelRepository, fuelCache)

	f := e.Group("/fuels")
	f.Use(middlewares.CheckAuth)
	f.GET("", fuelController.GetAllFuels)
	f.GET("/:id", fuelController.GetFuelByID)
	f.POST("", middlewares.CheckRoleAdmin(fuelController.CreateFuel))
	f.PUT("/:id", middlewares.CheckRoleAdmin(fuelController.UpdateFuel))
	f.DELETE("/:id", middlewares.CheckRoleAdmin(fuelController.DeleteFuel))

	// Electric routes
	electricRepository := repositories.NewElectricRepository(config.DB)
	electricCache := caches.NewElectricCache(config.RedisClient)
	electricController := controllers.NewElectricController(electricRepository, electricCache)

	l := e.Group("/electrics")
	l.Use(middlewares.CheckAuth)
	l.GET("", electricController.GetAllElectrics)
	l.GET("/:id", electricController.GetElectricByID)
	l.POST("", middlewares.CheckRoleAdmin(electricController.CreateElectric))
	l.PUT("/:id", middlewares.CheckRoleAdmin(electricController.UpdateElectric))
	l.DELETE("/:id", middlewares.CheckRoleAdmin(electricController.DeleteElectric))

	// Tree category routes
	ct := e.Group("/tree-categories")
	ct.Use(middlewares.CheckAuth)
	ct.GET("", controllers.GetAllTreeCategories)
	ct.GET("/:id", controllers.GetTreeCategoryByID)
	ct.POST("", middlewares.CheckRoleAdmin(controllers.CreateTreeCategory))
	ct.PUT("/:id", middlewares.CheckRoleAdmin(controllers.UpdateTreeCategory))
	ct.DELETE("/:id", middlewares.CheckRoleAdmin(controllers.DeleteTreeCategory))

	// Tree routes
	treeRepository := repositories.NewTreeRepository(config.DB)
	treeCache := caches.NewTreeCache(config.RedisClient)
	treeController := controllers.NewTreeController(treeRepository, treeCache)

	tr := e.Group("/trees")
	tr.Use(middlewares.CheckAuth)
	tr.GET("", treeController.GetAllTrees)
	tr.GET("/:id", treeController.GetTreeByID)
	tr.POST("", middlewares.CheckRoleAdmin(treeController.CreateTree))
	tr.PUT("/:id", middlewares.CheckRoleAdmin(treeController.UpdateTree))
	tr.DELETE("/:id", middlewares.CheckRoleAdmin(treeController.DeleteTree))

	// payment method
	paymentMethodRepository := repositories.NewPaymentMethodRepository(config.MongoCollection)
	paymentMethodController := controllers.NewPaymentMethodController(paymentMethodRepository)

	pm := e.Group("/payment-methods")
	pm.Use(middlewares.CheckAuth)
	pm.GET("", paymentMethodController.GetAllPaymentMethods)
	pm.POST("", middlewares.CheckRoleAdmin(paymentMethodController.CreatePaymentMethod))
	pm.PUT("/:id", middlewares.CheckRoleAdmin(paymentMethodController.UpdatePaymentMethod))
	pm.DELETE("/:id", middlewares.CheckRoleAdmin(paymentMethodController.DeletePaymentMethod))

	// Carbon fuel routes
	carbonFuelRepository := repositories.NewCarbonFuelRepository(config.DB)
	carbonFuelController := controllers.NewCarbonFuelController(carbonFuelRepository)

	cf := e.Group("/carbon-fuels")
	cf.Use(middlewares.CheckAuth)
	cf.Use(middlewares.CheckRoleCustomer)
	cf.GET("", carbonFuelController.GetAllCarbonFuels)
	cf.GET("/:id", carbonFuelController.GetCarbonFuelByID)
	cf.POST("", carbonFuelController.CreateCarbonFuel)
	cf.DELETE("/:id", carbonFuelController.DeleteCarbonFuel)

	// Carbon electric routes
	carbonElectricRepo := repositories.NewCarbonElectricRepository(config.DB)
	carbonElectricController := controllers.NewCarbonElectricController(carbonElectricRepo)

	ce := e.Group("/carbon-electrics")
	ce.Use(middlewares.CheckAuth)
	ce.GET("", carbonElectricController.GetAllCarbonElectrics)
	ce.GET("/:id", carbonElectricController.GetCarbonElectricByID)
	ce.POST("", carbonElectricController.CreateCarbonElectric)
	ce.DELETE("/:id", carbonElectricController.DeleteCarbonElectric)

	// Carbon summary routes
	carbonSummaryRepository := repositories.NewCarbonSummaryRepository(config.DB)
	carbonSummaryController := controllers.NewCarbonSummaryController(carbonSummaryRepository)

	cs := e.Group("/carbon-summaries")
	cs.Use(middlewares.CheckAuth)
	cs.Use(middlewares.CheckRoleCustomer)
	cs.GET("", carbonSummaryController.GetCarbonSummary)

	// Cart routes
	cartRepository := repositories.NewCartRepository(config.DB)
	cartController := controllers.NewCartController(cartRepository)

	c := e.Group("/carts")
	c.Use(middlewares.CheckAuth)
	c.Use(middlewares.CheckRoleCustomer)
	c.GET("", cartController.GetAllCart)
	c.POST("", cartController.AddCart)
	c.DELETE("/:id", cartController.DeleteCart)

	// Transaction routes
	transactionRepository := repositories.NewTransactionRepository(config.DB)
	transactionController := controllers.NewTransactionController(transactionRepository)

	t := e.Group("/transactions")
	t.Use(middlewares.CheckAuth)
	t.Use(middlewares.CheckRoleCustomer)
	t.GET("", transactionController.GetAllTransactions)
	t.POST("", transactionController.AddTransaction)

	// Payment routes
	paymentRepository := repositories.NewPaymentRepository(config.DB, config.MongoCollection)
	paymentController := controllers.NewPaymentController(paymentRepository)

	p := e.Group("/payments")
	p.POST("", middlewares.CheckAuth(middlewares.CheckRoleCustomer(paymentController.CreatePayment)))
	p.GET("/verify/:id", paymentController.VerifyPayment)

	// Report routes
	reportRepository := repositories.NewReportRepository(config.DB)
	reportController := controllers.NewReportController(reportRepository)

	rp := e.Group("/reports")
	rp.Use(middlewares.CheckAuth)
	rp.GET("/summary", reportController.GetReportSummary)

	// Generate PDF routes
	pdfGeneratorController := controllers.NewGeneratePdfController()

	e.POST("/generate-pdf", pdfGeneratorController.PdfHandler)
	e.POST("/generate-pdf-summary", pdfGeneratorController.PdfHandlerSummary)

	// Gemini API routes
	newGeminiAPIController := controllers.NewGeminiAPIController()

	e.POST("/ai", newGeminiAPIController.GeminiAPI)
	e.POST("/ai/generate-image", newGeminiAPIController.GenerateImage)

	renderer := controllers.NewTemplateRenderer("views")
	e.Renderer = renderer
	emailVerifyController := controllers.NewEmailVerificationController(userRepository)
	e.Debug = true
	e.GET("/verify-email", emailVerifyController.HandleEmailVerification)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
