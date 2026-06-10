package main

import (
	"cinema-backend/internal/config"
	"cinema-backend/internal/handlers"
	"cinema-backend/internal/middleware"
	"cinema-backend/internal/repository"
	"cinema-backend/internal/services"
	ws "cinema-backend/internal/websocket"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := repository.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	theaterRepo := repository.NewTheaterRepository(db)
	screenRepo := repository.NewScreenRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	showtimeRepo := repository.NewShowtimeRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	concessionRepo := repository.NewConcessionRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	seatLockRepo := repository.NewSeatLockRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	movieService := services.NewMovieService(movieRepo)
	screenService := services.NewScreenService(screenRepo, seatRepo)
	theaterService := services.NewTheaterService(theaterRepo, screenService)
	seatService := services.NewSeatService(seatRepo, bookingRepo, seatLockRepo)
	showtimeService := services.NewShowtimeService(showtimeRepo)
	ticketService := services.NewTicketService(ticketRepo, showtimeRepo)
	concessionService := services.NewConcessionService(concessionRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	emailService := services.NewEmailService()
	bookingService := services.NewBookingService(bookingRepo, showtimeRepo, seatLockRepo, paymentRepo, ticketRepo, emailService)
	vnpayService := services.NewVNPayService(cfg.VNPay)
	seatLockService := services.NewSeatLockService(seatLockRepo, bookingRepo)

	cloudinaryService, err := services.NewCloudinaryService(cfg.CloudinaryURL)
	if err != nil {
		log.Printf("Warning: failed to initialize Cloudinary: %v", err)
	}
	if cloudinaryService != nil && cloudinaryService.IsEnabled() {
		log.Println("Cloudinary integration enabled")
	} else {
		log.Println("Cloudinary integration disabled - using local file storage")
	}

	// Start cleanup routine for expired seat locks
	seatLockService.StartCleanupRoutine()

	// Start cleanup routines (every hour)
	go func() {
		for {
			showtimeService.DeleteExpired()

			retentionDays := 7
			setting, err := settingRepo.GetByKey("ticket_retention_days")
			if err == nil {
				if d, err := strconv.Atoi(setting.Value); err == nil && d > 0 {
					retentionDays = d
				}
			}
			if err := ticketRepo.DeleteExpired(retentionDays); err != nil {
				log.Printf("Error deleting expired tickets: %v", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	// Initialize WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()
	wsHandler := ws.NewWebSocketHandler(hub, seatService)

	// Set WebSocket Hub for seat lock service to enable real-time updates
	seatLockService.SetHub(hub)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	movieHandler := handlers.NewMovieHandler(movieService)
	theaterHandler := handlers.NewTheaterHandler(theaterService)
	screenHandler := handlers.NewScreenHandler(screenService)
	seatHandler := handlers.NewSeatHandler(seatService)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	concessionHandler := handlers.NewConcessionHandler(concessionService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	paymentHandler := handlers.NewPaymentHandler(vnpayService, bookingService, paymentService)
	seatLockHandler := handlers.NewSeatLockHandler(seatLockService)
	baseURL := os.Getenv("RENDER_EXTERNAL_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + cfg.Port
	}
	uploadHandler := handlers.NewUploadHandler(baseURL+"/api", cloudinaryService)
	settingHandler := handlers.NewSettingHandler(settingRepo)

	// Set DB for dashboard handlers
	handlers.SetDashboardDB(db)

	// Setup router with CORS first
	r := gin.Default()

	// CORS middleware - MUST be first
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint (no auth required)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public movie routes
		api.GET("/movies", movieHandler.GetAll)
		api.GET("/movies/now-showing", movieHandler.GetNowShowing)
		api.GET("/movies/coming-soon", movieHandler.GetComingSoon)
		api.GET("/movies/:id", movieHandler.GetByID)

		// Public theater routes
		api.GET("/theaters", theaterHandler.GetAll)
		api.GET("/theaters/:id", theaterHandler.GetByID)
		api.GET("/screens/theater/:theaterId", screenHandler.GetByTheater)

		// Public showtime routes
		api.GET("/showtimes/movie/:movieId", showtimeHandler.GetByMovie)
		api.GET("/showtimes/movie/:movieId/theater/:theaterId", showtimeHandler.GetByMovieAndTheater)
		api.GET("/showtimes/:id", showtimeHandler.GetByID)

		// Public seat routes
		api.GET("/seats/screen/:screenId", seatHandler.GetByScreen)

		// WebSocket endpoint
		api.GET("/ws", wsHandler.HandleWebSocket)

		// Public concession routes
		api.GET("/concessions", concessionHandler.GetAll)
		api.GET("/concessions/category/:category", concessionHandler.GetByCategory)

		// VNPay payment callbacks (public)
		api.GET("/payments/vnpay/return", paymentHandler.VNPayReturn)

		// Mock payment pages for sandbox testing
		api.GET("/payments/vnpay/mock", paymentHandler.MockPaymentPage)
		api.GET("/payments/vnpay/mock-success", paymentHandler.MockPaymentSuccess)
		api.GET("/payments/vnpay/mock-fail", paymentHandler.MockPaymentFail)

		// Public image serving
		api.GET("/uploads/:filename", uploadHandler.ServeImage)
		api.HEAD("/uploads/:filename", uploadHandler.ServeImage)
	}

	// Protected routes
	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Auth
		authorized.GET("/auth/me", authHandler.Me)
		authorized.POST("/auth/logout", authHandler.Logout)
		authorized.PUT("/auth/profile", authHandler.UpdateProfile)
		authorized.POST("/auth/change-password", authHandler.ChangePassword)

		// Admin only routes
		admin := authorized.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			// Movies
			admin.POST("/movies", movieHandler.Create)
			admin.PUT("/movies/:id", movieHandler.Update)
			admin.DELETE("/movies/:id", movieHandler.Delete)

			// Upload
			admin.POST("/uploads", uploadHandler.UploadImage)
			admin.DELETE("/uploads/:filename", uploadHandler.DeleteImage)

			// Theaters (admin only - create/update/delete)
			admin.POST("/theaters", theaterHandler.Create)
			admin.PUT("/theaters/:id", theaterHandler.Update)
			admin.DELETE("/theaters/:id", theaterHandler.Delete)

			// Screens (admin only - create/update/delete)
			admin.GET("/screens", screenHandler.GetAll)
			admin.POST("/screens", screenHandler.Create)
			admin.GET("/screens/:id", screenHandler.GetByID)
			admin.PUT("/screens/:id", screenHandler.Update)
			admin.DELETE("/screens/:id", screenHandler.Delete)

			// Concessions (admin only - create/update/delete)
			admin.POST("/concessions", concessionHandler.Create)
			admin.PUT("/concessions/:id", concessionHandler.Update)
			admin.DELETE("/concessions/:id", concessionHandler.Delete)

			// Showtimes
			admin.GET("/showtimes", showtimeHandler.GetAll)
			admin.POST("/showtimes", showtimeHandler.Create)
			admin.PUT("/showtimes/:id", showtimeHandler.Update)
			admin.DELETE("/showtimes/:id", showtimeHandler.Delete)

			// Bookings (admin view)
			admin.GET("/bookings", bookingHandler.GetAll)
			admin.GET("/bookings/:id", bookingHandler.GetByID)
			admin.POST("/bookings/:id/confirm", bookingHandler.Confirm)
			admin.POST("/bookings/:id/cancel", bookingHandler.Cancel)
			admin.POST("/bookings/:id/refund", bookingHandler.Refund)

			// Users
			admin.GET("/users", authHandler.GetAllUsers)
			admin.GET("/users/customers", authHandler.GetCustomers)
			admin.GET("/users/:id", authHandler.GetUserByID)
			admin.GET("/users/:id/bookings", bookingHandler.GetMyBookings)
			admin.DELETE("/users/:id", authHandler.DeleteUser)

			// Settings
			admin.GET("/settings", settingHandler.GetAll)
			admin.PUT("/settings", settingHandler.Update)

			// Dashboard
			admin.GET("/dashboard/admin-stats", handlers.GetAdminStats)
			admin.GET("/dashboard/recent-bookings", handlers.GetRecentBookings)
			admin.GET("/dashboard/top-movies", handlers.GetTopMovies)
			admin.GET("/dashboard/revenue-by-day", handlers.GetRevenueByDay)
			admin.GET("/dashboard/revenue-by-movie", handlers.GetRevenueByMovie)
			admin.GET("/dashboard/revenue-by-theater", handlers.GetRevenueByTheater)
			admin.GET("/dashboard/revenue-by-genre", handlers.GetRevenueByGenre)
			admin.GET("/dashboard/bookings-by-genre", handlers.GetBookingsByGenre)
		}

		// Customer routes (also accessible by admin)
		customer := authorized.Group("/")
		customer.Use(middleware.CustomerMiddleware())
		{
			// Tickets
			customer.GET("/tickets/my-tickets", ticketHandler.GetMyTickets)
			customer.GET("/tickets/:id", ticketHandler.GetByID)
			customer.POST("/tickets/book", ticketHandler.Book)
			customer.POST("/tickets/:id/cancel", ticketHandler.Cancel)

			// Bookings
			customer.GET("/my-bookings", bookingHandler.GetMyBookings)
			customer.POST("/bookings/create", bookingHandler.Create)
			customer.DELETE("/my-bookings/clear", bookingHandler.ClearMyBookings)

			// Seat Locks
			customer.POST("/seat-locks", seatLockHandler.LockSeats)
			customer.DELETE("/seat-locks/:id", seatLockHandler.UnlockSeats)
			customer.PUT("/seat-locks/:id/extend", seatLockHandler.ExtendLock)
			customer.GET("/seat-locks/:id", seatLockHandler.GetLockStatus)
			customer.GET("/seat-locks", seatLockHandler.GetActiveLocks)
			customer.POST("/seat-locks/clear-all", seatLockHandler.ClearAllLocks)

			// Dashboard
			customer.GET("/dashboard/customer-stats", handlers.GetCustomerStats)

			// VNPay Payments
			customer.POST("/payments/vnpay/create", paymentHandler.CreateVNPayPayment)
			customer.GET("/payments/status/:bookingId", paymentHandler.CheckPaymentStatus)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
