package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BookingService struct {
	repo         *repository.BookingRepository
	showtimeRepo *repository.ShowtimeRepository
	lockRepo     *repository.SeatLockRepository
	paymentRepo  *repository.PaymentRepository
	ticketRepo   *repository.TicketRepository
	emailService *EmailService
}

func NewBookingService(repo *repository.BookingRepository, showtimeRepo *repository.ShowtimeRepository, lockRepo *repository.SeatLockRepository, paymentRepo *repository.PaymentRepository, ticketRepo *repository.TicketRepository, emailService *EmailService) *BookingService {
	return &BookingService{
		repo:         repo,
		showtimeRepo: showtimeRepo,
		lockRepo:     lockRepo,
		paymentRepo:  paymentRepo,
		ticketRepo:   ticketRepo,
		emailService: emailService,
	}
}

func (s *BookingService) Create(booking *models.Booking, userID uuid.UUID) error {
	// Validate showtime exists
	_, err := s.showtimeRepo.FindByID(booking.ShowtimeID)
	if err != nil {
		return err
	}

	// Check if seats are still available (not booked by others)
	bookedSeatIDs, err := s.repo.GetBookedSeatsByShowtime(booking.ShowtimeID)
	if err != nil {
		return err
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, id := range bookedSeatIDs {
		bookedMap[id] = true
	}

	for _, seat := range booking.BookingSeats {
		if bookedMap[seat.SeatID] {
			return errors.New("one or more seats are already booked")
		}
	}

	// Check if user has an active lock for these seats
	lock, err := s.lockRepo.FindActiveByUser(userID, booking.ShowtimeID)
	if err == nil && lock != nil {
		// Verify the locked seats match the booking seats
		lockSeatMap := make(map[uuid.UUID]bool)
		for _, id := range lock.SeatIDs {
			lockSeatMap[id] = true
		}

		for _, seat := range booking.BookingSeats {
			if !lockSeatMap[seat.SeatID] {
				return errors.New("seats do not match your reservation")
			}
		}

		// Release the lock after booking
		defer s.lockRepo.ReleaseLock(lock.ID)
	} else {
		// If no lock exists, check if seats are locked by others
		for _, seat := range booking.BookingSeats {
			existingLock, err := s.lockRepo.FindBySeatAndShowtime(seat.SeatID, booking.ShowtimeID)
			if err == nil && existingLock != nil {
				return errors.New("one or more seats are reserved by another user")
			}
		}
	}

	// Set default status to PENDING (đã giữ ghế, chờ thanh toán)
	if booking.Status == "" {
		booking.Status = models.BookingPending
	}

	// Set default payment status
	if booking.PaymentStatus == "" {
		booking.PaymentStatus = models.PaymentPending
	}

	// Generate booking code
	if booking.BookingCode == "" {
		booking.BookingCode = generateBookingCode()
	}

	return s.repo.Create(booking)
}

func (s *BookingService) GetAll(page, limit int) ([]models.Booking, int64, error) {
	return s.repo.FindAll(page, limit)
}

func (s *BookingService) GetMyBookings(userID uuid.UUID) ([]models.Booking, error) {
	return s.repo.FindByUser(userID)
}

func (s *BookingService) GetByID(id string) (*models.Booking, error) {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(bookingID)
}

// ConfirmPayment xử lý khi thanh toán thành công
func (s *BookingService) ConfirmPayment(id string, paymentInfo map[string]string) error {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Update booking status
	if err := s.repo.UpdateStatus(bookingID, models.BookingConfirmed); err != nil {
		return err
	}

	// Update payment status
	if err := s.repo.UpdatePaymentStatus(bookingID, models.PaymentPaid); err != nil {
		return err
	}

	// Create ticket
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return err
	}

	// Generate QR code
	qrCode := generateQRCode(booking.BookingCode)

	// Update booking with QR
	booking.QRCode = &qrCode
	if err := s.repo.Update(booking); err != nil {
		return err
	}

	// Create ticket record
	ticket := &models.Ticket{
		BookingID:   &bookingID,
		ShowtimeID:  booking.ShowtimeID,
		UserID:      booking.UserID,
		TotalPrice:  booking.TotalAmount,
		Status:      models.TicketPaid,
		QRCode:      &qrCode,
		BookingTime: time.Now(),
	}

	// Convert seats to JSON
	seatsJSON := convertSeatsToJSON(booking.BookingSeats)
	ticket.Seats = seatsJSON

	if err := s.ticketRepo.Create(ticket); err != nil {
		return err
	}

	// Send email confirmation
	go s.sendTicketConfirmationEmail(booking, qrCode)

	return nil
}

// sendTicketConfirmationEmail gửi email xác nhận vé
func (s *BookingService) sendTicketConfirmationEmail(booking *models.Booking, qrCode string) {
	if s.emailService == nil {
		return
	}

	// Get user email
	user, err := s.repo.FindByID(booking.ID)
	if err != nil {
		fmt.Printf("Error finding user for email: %v\n", err)
		return
	}

	// Format showtime
	showtimeStr := "N/A"
	if booking.Showtime.StartTime.Year() > 1 {
		showtimeStr = booking.Showtime.StartTime.Format("2006-01-02 15:04")
	}

	// Format seats
	seats := ""
	for i, seat := range booking.BookingSeats {
		if i > 0 {
			seats += ", "
		}
		seats += seat.SeatLabel
	}

	bookingDetails := map[string]interface{}{
		"bookingCode": booking.BookingCode,
		"movie":       booking.Showtime.Movie.Title,
		"showtime":    showtimeStr,
		"seats":       seats,
		"amount":      booking.TotalAmount,
		"qrCode":      qrCode,
	}

	if err := s.emailService.SendTicketConfirmation(user.User.Email, user.User.Name, bookingDetails); err != nil {
		fmt.Printf("Error sending ticket confirmation email: %v\n", err)
	} else {
		fmt.Printf("Ticket confirmation email sent to %s\n", user.User.Email)
	}
}

// CancelBooking hủy booking
func (s *BookingService) CancelBooking(id string, reason string) error {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Update booking status
	if err := s.repo.UpdateStatus(bookingID, models.BookingCancelled); err != nil {
		return err
	}

	// Update payment status to refunded if was paid
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return err
	}

	if booking.PaymentStatus == models.PaymentPaid {
		if err := s.repo.UpdatePaymentStatus(bookingID, models.PaymentRefunded); err != nil {
			return err
		}
	}

	return nil
}

// ExpireBooking xử lý khi booking hết hạn thanh toán
func (s *BookingService) ExpireBooking(id string) error {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Update booking status to expired
	if err := s.repo.UpdateStatus(bookingID, models.BookingExpired); err != nil {
		return err
	}

	// Update payment status to failed
	if err := s.repo.UpdatePaymentStatus(bookingID, models.PaymentFailed); err != nil {
		return err
	}

	return nil
}

// CompleteBooking xử lý khi suất chiếu kết thúc
func (s *BookingService) CompleteBooking(id string) error {
	bookingID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(bookingID, models.BookingCompleted)
}

// ProcessExpiredBookings xử lý các booking đã hết hạn
func (s *BookingService) ProcessExpiredBookings() error {
	// Tìm các booking PENDING đã quá 15 phút
	expirationTime := time.Now().Add(-15 * time.Minute)

	var expiredBookings []models.Booking
	if err := s.repo.FindPendingExpired(expirationTime, &expiredBookings); err != nil {
		return err
	}

	for _, booking := range expiredBookings {
		if err := s.ExpireBooking(booking.ID.String()); err != nil {
			// Log error but continue processing others
			fmt.Printf("Error expiring booking %s: %v\n", booking.ID, err)
		}
	}

	return nil
}

// Refund xử lý hoàn tiền
func (s *BookingService) Refund(id string) error {
	return s.CancelBooking(id, "refunded")
}

func (s *BookingService) DeleteMyBookings(userID uuid.UUID) error {
	return s.repo.DeleteByUser(userID)
}

// Helper functions
func generateBookingCode() string {
	return fmt.Sprintf("BK%s", time.Now().Format("20060102150405"))
}

func generateQRCode(bookingCode string) string {
	return fmt.Sprintf("QR-%s-%d", bookingCode, time.Now().Unix())
}

func convertSeatsToJSON(seats []models.BookingSeat) string {
	// Simple JSON representation
	result := "["
	for i, seat := range seats {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`{"row":"%s","number":%d,"price":%.2f}`, seat.SeatLabel, i+1, seat.Price)
	}
	result += "]"
	return result
}
