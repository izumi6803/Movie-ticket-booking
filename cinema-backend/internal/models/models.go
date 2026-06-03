package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"default:'customer'"`
	Phone     *string        `json:"phone,omitempty"`
	Avatar    *string        `json:"avatar,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Tickets  []Ticket  `json:"tickets,omitempty" gorm:"foreignKey:UserID"`
	Bookings []Booking `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
}

type MovieGenre string

const (
	GenreAction      MovieGenre = "action"
	GenreAdventure   MovieGenre = "adventure"
	GenreAnimation   MovieGenre = "animation"
	GenreComedy      MovieGenre = "comedy"
	GenreCrime       MovieGenre = "crime"
	GenreDocumentary MovieGenre = "documentary"
	GenreDrama       MovieGenre = "drama"
	GenreFamily      MovieGenre = "family"
	GenreFantasy     MovieGenre = "fantasy"
	GenreHorror      MovieGenre = "horror"
	GenreMusical     MovieGenre = "musical"
	GenreMystery     MovieGenre = "mystery"
	GenreRomance     MovieGenre = "romance"
	GenreSciFi       MovieGenre = "sci-fi"
	GenreThriller    MovieGenre = "thriller"
	GenreWar         MovieGenre = "war"
	GenreWestern     MovieGenre = "western"
)

type MovieRating string

const (
	RatingG    MovieRating = "G"
	RatingPG   MovieRating = "PG"
	RatingPG13 MovieRating = "PG-13"
	RatingR    MovieRating = "R"
	RatingNC17 MovieRating = "NC-17"
)

type MovieStatus string

const (
	MovieNowShowing MovieStatus = "now_showing"
	MovieComingSoon MovieStatus = "coming_soon"
	MovieEnded      MovieStatus = "ended"
)

type Movie struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Duration    int            `json:"duration" gorm:"not null"` // in minutes
	Genre       []string       `json:"genre,omitempty" gorm:"type:jsonb;serializer:json"`
	Rating      MovieRating    `json:"rating,omitempty"`
	PosterURL   *string        `json:"posterUrl,omitempty"`
	TrailerURL  *string        `json:"trailerUrl,omitempty"`
	ReleaseDate *time.Time     `json:"releaseDate,omitempty"`
	Director    string         `json:"director,omitempty"`
	Cast        []string       `json:"cast,omitempty" gorm:"type:jsonb;serializer:json"`
	Status      MovieStatus    `json:"status" gorm:"default:'coming_soon'"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Showtimes []Showtime `json:"showtimes,omitempty" gorm:"foreignKey:MovieID"`
}

type ScreenType string

const (
	ScreenStandard ScreenType = "standard"
	ScreenIMAX     ScreenType = "imax"
	Screen3D       ScreenType = "3d"
	Screen4DX      ScreenType = "4dx"
	ScreenVIP      ScreenType = "vip"
)

type HallStatus string

const (
	HallActive      HallStatus = "active"
	HallMaintenance HallStatus = "maintenance"
	HallInactive    HallStatus = "inactive"
)

type CinemaHall struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	TotalSeats  int            `json:"totalSeats" gorm:"not null"`
	RowsCount   int            `json:"rows" gorm:"column:rows_count;not null"`
	SeatsPerRow int            `json:"seatsPerRow" gorm:"not null"`
	ScreenType  ScreenType     `json:"screenType" gorm:"default:'standard'"`
	SoundSystem string         `json:"soundSystem"`
	Status      HallStatus     `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Showtimes []Showtime `json:"showtimes,omitempty" gorm:"foreignKey:HallID"`
}

type ShowtimeStatus string

const (
	ShowtimeScheduled ShowtimeStatus = "scheduled"
	ShowtimeOngoing   ShowtimeStatus = "ongoing"
	ShowtimeCompleted ShowtimeStatus = "completed"
	ShowtimeCancelled ShowtimeStatus = "cancelled"
)

type Showtime struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MovieID         uuid.UUID      `json:"movieId" gorm:"not null"`
	ScreenID        uuid.UUID      `json:"screenId" gorm:"not null"`
	StartTime       time.Time      `json:"startTime" gorm:"not null"`
	EndTime         time.Time      `json:"endTime" gorm:"not null"`
	BaseTicketPrice float64        `json:"baseTicketPrice" gorm:"not null"`
	AvailableSeats  int            `json:"availableSeats" gorm:"not null"`
	Status          ShowtimeStatus `json:"status" gorm:"default:'active'"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Movie   Movie    `json:"movie,omitempty" gorm:"foreignKey:MovieID"`
	Screen  Screen   `json:"screen,omitempty" gorm:"foreignKey:ScreenID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:ShowtimeID"`
}

type SeatType string

const (
	SeatStandard   SeatType = "standard"
	SeatPremium    SeatType = "premium"
	SeatVIP        SeatType = "vip"
	SeatWheelchair SeatType = "wheelchair"
)

type TicketStatus string

const (
	TicketPending   TicketStatus = "pending"
	TicketPaid      TicketStatus = "paid"
	TicketCancelled TicketStatus = "cancelled"
	TicketRefunded  TicketStatus = "refunded"
	TicketUsed      TicketStatus = "used"
)

type Ticket struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID   *uuid.UUID     `json:"bookingId,omitempty" gorm:"index"`
	ShowtimeID  uuid.UUID      `json:"showtimeId" gorm:"not null"`
	UserID      uuid.UUID      `json:"userId" gorm:"not null"`
	Seats       string         `json:"seats" gorm:"type:jsonb;not null"` // JSON array of seat objects
	TotalPrice  float64        `json:"totalPrice" gorm:"not null"`
	Status      TicketStatus   `json:"status" gorm:"default:'pending'"`
	QRCode      *string        `json:"qrCode,omitempty"`
	BookingTime time.Time      `json:"bookingTime"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Showtime Showtime `json:"showtime,omitempty" gorm:"foreignKey:ShowtimeID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Booking  *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

type Theater struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string         `json:"name" gorm:"not null"`
	Location     string         `json:"location"`
	TotalScreens int            `json:"totalScreens" gorm:"default:1"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Screens []Screen `json:"screens,omitempty" gorm:"foreignKey:TheaterID"`
}

type Screen struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TheaterID   uuid.UUID      `json:"theaterId" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	ScreenType  ScreenType     `json:"screenType" gorm:"default:'standard'"`
	TotalRows   int            `json:"totalRows" gorm:"not null"`
	SeatsPerRow int            `json:"seatsPerRow" gorm:"not null"`
	TotalSeats  int            `json:"totalSeats" gorm:"not null"`
	SoundSystem string         `json:"soundSystem"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Theater   Theater    `json:"theater,omitempty" gorm:"foreignKey:TheaterID"`
	Showtimes []Showtime `json:"showtimes,omitempty" gorm:"foreignKey:ScreenID"`
	Seats     []Seat     `json:"seats,omitempty" gorm:"foreignKey:ScreenID"`
}

type Seat struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScreenID        uuid.UUID      `json:"screenId" gorm:"not null"`
	RowLabel        string         `json:"rowLabel" gorm:"not null"`
	SeatNumber      int            `json:"seatNumber" gorm:"not null"`
	SeatType        SeatType       `json:"seatType" gorm:"default:'standard'"`
	PriceMultiplier float64        `json:"priceMultiplier" gorm:"default:1.00"`
	CreatedAt       time.Time      `json:"createdAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type ConcessionCategory string

const (
	ConcessionFood        ConcessionCategory = "food"
	ConcessionDrink       ConcessionCategory = "drink"
	ConcessionMerchandise ConcessionCategory = "merchandise"
	ConcessionCombo       ConcessionCategory = "combo"
)

type Concession struct {
	ID            uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string             `json:"name" gorm:"not null"`
	Description   string             `json:"description"`
	Category      ConcessionCategory `json:"category" gorm:"not null"`
	Price         float64            `json:"price" gorm:"not null"`
	ImageURL      *string            `json:"imageUrl,omitempty"`
	StockQuantity int                `json:"stockQuantity" gorm:"default:0"`
	IsActive      bool               `json:"isActive" gorm:"default:true"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt     `json:"-" gorm:"index"`
}

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"   // Đã giữ ghế, chờ thanh toán
	BookingConfirmed BookingStatus = "confirmed" // Thanh toán thành công, vé hợp lệ
	BookingCancelled BookingStatus = "cancelled" // User/admin hủy
	BookingExpired   BookingStatus = "expired"   // Quá thởi gian thanh toán, ghế được mở lại
	BookingCompleted BookingStatus = "completed" // Suất chiếu đã kết thúc
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"  // Chờ thanh toán
	PaymentPaid     PaymentStatus = "paid"     // Thanh toán thành công
	PaymentFailed   PaymentStatus = "failed"   // Thanh toán thất bại
	PaymentRefunded PaymentStatus = "refunded" // Đã hoàn tiền
)

type Booking struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID               uuid.UUID      `json:"userId" gorm:"not null"`
	ShowtimeID           uuid.UUID      `json:"showtimeId" gorm:"not null"`
	BookingCode          string         `json:"bookingCode" gorm:"uniqueIndex;not null"`
	TotalTicketPrice     float64        `json:"totalTicketPrice" gorm:"not null;default:0"`
	TotalConcessionPrice float64        `json:"totalConcessionPrice" gorm:"not null;default:0"`
	TotalAmount          float64        `json:"totalAmount" gorm:"not null;default:0"`
	Status               BookingStatus  `json:"status" gorm:"default:'pending'"`
	PaymentStatus        PaymentStatus  `json:"paymentStatus" gorm:"default:'pending'"`
	QRCode               *string        `json:"qrCode,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Showtime     Showtime      `json:"showtime,omitempty" gorm:"foreignKey:ShowtimeID"`
	BookingSeats []BookingSeat `json:"bookingSeats,omitempty" gorm:"foreignKey:BookingID"`
	OrderItems   []OrderItem   `json:"orderItems,omitempty" gorm:"foreignKey:BookingID"`
}

type BookingSeat struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID uuid.UUID `json:"bookingId" gorm:"not null"`
	SeatID    uuid.UUID `json:"seatId" gorm:"not null"`
	SeatLabel string    `json:"seatLabel" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`

	// Relations
	Booking Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	Seat    Seat    `json:"seat,omitempty" gorm:"foreignKey:SeatID"`
}

type OrderItem struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID    uuid.UUID `json:"bookingId" gorm:"not null"`
	ConcessionID uuid.UUID `json:"concessionId" gorm:"not null"`
	Quantity     int       `json:"quantity" gorm:"not null;default:1"`
	UnitPrice    float64   `json:"unitPrice" gorm:"not null"`
	TotalPrice   float64   `json:"totalPrice" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`

	// Relations
	Booking    Booking    `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	Concession Concession `json:"concession,omitempty" gorm:"foreignKey:ConcessionID"`
}
