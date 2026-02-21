package main

import (
	"log"
	"net/http"
	"time"

	"arem-shop/internal/config"
	"arem-shop/internal/database"
	"arem-shop/internal/handlers"
	"arem-shop/internal/middleware"
	"arem-shop/internal/models"
	"arem-shop/internal/repository"
	"arem-shop/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	shopRepo := repository.NewShopRepository(db)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	authService := services.NewAuthService(cfg, userRepo, shopRepo)
	productService := services.NewProductService(productRepo)
	publicService := services.NewPublicService(shopRepo, productRepo)
	transactionService := services.NewTransactionService(db, productRepo, transactionRepo)
	reportService := services.NewReportService(transactionRepo, productRepo, cfg.LowStockThreshold)
	categoryService := services.NewCategoryService(categoryRepo)
	shopService := services.NewShopService(shopRepo)

	// URL de base exposée dans l'API, utilisée pour construire les URLs d'images retournées
	baseURL := "http://localhost:8080"
	uploadService := services.NewUploadService("./uploads", baseURL)

	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	publicHandler := handlers.NewPublicHandler(publicService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	shopHandler := handlers.NewShopHandler(shopService)

	router.GET("/health", func(c *gin.Context) {
		sqlDB, dbErr := db.DB()
		dbStatus := "up"
		if dbErr != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}

		c.JSON(http.StatusOK, gin.H{
			"app":         cfg.AppName,
			"environment": cfg.AppEnv,
			"database":    dbStatus,
			"timestamp":   time.Now().UTC(),
		})
	})

	auth := router.Group("/auth")
	auth.POST("/login", authHandler.Login)

	protectedAuth := router.Group("/auth")
	protectedAuth.Use(
		middleware.AuthMiddleware(cfg),
		middleware.ShopIsolationMiddleware(),
		middleware.RoleMiddleware(models.RoleSuperAdmin),
	)
	protectedAuth.POST("/register", authHandler.Register)

	private := router.Group("/")
	private.Use(
		middleware.AuthMiddleware(cfg),
		middleware.ShopIsolationMiddleware(),
		middleware.RoleMiddleware(models.RoleSuperAdmin, models.RoleAdmin),
	)
	private.GET("/categories", categoryHandler.List)
	private.GET("/products", productHandler.List)
	private.GET("/products/:id", productHandler.GetByID)
	private.POST("/products", productHandler.Create)
	private.PUT("/products/:id", productHandler.Update)
	private.DELETE("/products/:id", productHandler.Delete)
	private.POST("/transactions", transactionHandler.Create)
	private.POST("/upload", uploadHandler.UploadProductImage)

	// Servir le dossier './uploads' publiquement via l'URL '/uploads'
	router.Static("/uploads", "./uploads")

	superAdminPrivate := router.Group("/")
	superAdminPrivate.Use(
		middleware.AuthMiddleware(cfg),
		middleware.ShopIsolationMiddleware(),
		middleware.RoleMiddleware(models.RoleSuperAdmin),
	)
	superAdminPrivate.GET("/reports/dashboard", reportHandler.Dashboard)
	superAdminPrivate.PUT("/shop", shopHandler.UpdateShopInfo)

	public := router.Group("/public")
	public.GET("/:shopID/products", publicHandler.ListPublicProducts)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
