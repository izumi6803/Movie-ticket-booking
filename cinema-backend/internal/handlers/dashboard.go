package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Dashboard handlers with real database queries

var dashboardDB *gorm.DB

func SetDashboardDB(db *gorm.DB) {
	dashboardDB = db
}

func GetAdminStats(c *gin.Context) {
	if dashboardDB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "database not initialized"})
		return
	}

	var totalMovies int64
	dashboardDB.Table("movies").Count(&totalMovies)

	var nowShowing int64
	dashboardDB.Table("movies").Where("status = ?", "now_showing").Count(&nowShowing)

	var totalBookings int64
	dashboardDB.Table("bookings").Count(&totalBookings)

	var todayBookings int64
	dashboardDB.Raw("SELECT COUNT(*) FROM bookings WHERE DATE(created_at) = CURRENT_DATE").Scan(&todayBookings)

	var totalRevenue float64
	dashboardDB.Raw("SELECT COALESCE(SUM(total_amount), 0) FROM bookings WHERE status = ?", "paid").Scan(&totalRevenue)

	var todayRevenue float64
	dashboardDB.Raw("SELECT COALESCE(SUM(total_amount), 0) FROM bookings WHERE DATE(created_at) = CURRENT_DATE AND status = ?", "paid").Scan(&todayRevenue)

	var totalCustomers int64
	dashboardDB.Table("users").Where("role = ?", "customer").Count(&totalCustomers)

	stats := map[string]interface{}{
		"totalMovies":    totalMovies,
		"nowShowing":     nowShowing,
		"totalBookings":  totalBookings,
		"todayBookings":  todayBookings,
		"totalRevenue":   totalRevenue,
		"todayRevenue":   todayRevenue,
		"totalCustomers": totalCustomers,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func GetCustomerStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
		return
	}

	var totalBookings int64
	var totalSpent float64

	if dashboardDB != nil {
		dashboardDB.Table("bookings").Where("user_id = ?", userID).Count(&totalBookings)
		dashboardDB.Raw("SELECT COALESCE(SUM(total_amount), 0) FROM bookings WHERE user_id = ? AND status = ?", userID, "paid").Scan(&totalSpent)
	}

	stats := map[string]interface{}{
		"totalBookings": totalBookings,
		"upcomingShows": 0,
		"totalSpent":    totalSpent,
		"loyaltyPoints": 0,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func GetRecentBookings(c *gin.Context) {
	if dashboardDB == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
		return
	}

	var bookings []map[string]interface{}
	dashboardDB.Raw(`
		SELECT b.id, b.booking_code, b.status, b.total_amount, b.created_at,
			u.name as user_name, u.email as user_email,
			m.title as movie_title
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN showtimes s ON b.showtime_id = s.id
		JOIN movies m ON s.movie_id = m.id
		ORDER BY b.created_at DESC
		LIMIT 10
	`).Scan(&bookings)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": bookings})
}

func GetTopMovies(c *gin.Context) {
	if dashboardDB == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
		return
	}

	var movies []map[string]interface{}
	dashboardDB.Raw(`
		SELECT m.title as movie_title, COUNT(b.id) as bookings
		FROM bookings b
		JOIN showtimes s ON b.showtime_id = s.id
		JOIN movies m ON s.movie_id = m.id
		WHERE b.status = 'paid'
		GROUP BY m.id, m.title
		ORDER BY bookings DESC
		LIMIT 5
	`).Scan(&movies)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

func GetRevenueByDay(c *gin.Context) {
	if dashboardDB == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
		return
	}

	var revenue []map[string]interface{}
	dashboardDB.Raw(`
		SELECT DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as revenue
		FROM bookings
		WHERE status = 'paid'
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`).Scan(&revenue)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": revenue})
}

func GetBookingsByGenre(c *gin.Context) {
	if dashboardDB == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
		return
	}

	var genres []map[string]interface{}
	dashboardDB.Raw(`
		SELECT genre, COUNT(*) as bookings
		FROM bookings b
		JOIN showtimes s ON b.showtime_id = s.id
		JOIN movies m ON s.movie_id = m.id
		CROSS JOIN LATERAL UNNEST(m.genre) as genre
		WHERE b.status = 'paid'
		GROUP BY genre
		ORDER BY bookings DESC
	`).Scan(&genres)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": genres})
}
