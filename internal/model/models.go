package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"size:100;not null;unique" json:"email"`
	FullName    string `gorm:"size:100;not null" json:"full_name"`
	Password    string `gorm:"size:255;not null" json:"password"`
	PhoneNumber string `gorm:"type:text;not null;unique" json:"phone_number"`

	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

// Role model: bảng Roles
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Quan hệ many2many với Users và Permissions
	Users       []User       `gorm:"many2many:user_roles;" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

// Permission model: bảng Permissions
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Quan hệ many2many với Roles
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles"`
}

// RolePermission join table: bảng RolePermissions
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;not null" json:"role_id"`
	PermissionID uint `gorm:"primaryKey;not null" json:"permission_id"`
}

// UserRole join table: bảng UserRoles
type UserRole struct {
	UserID uint `gorm:"primaryKey;not null" json:"user_id"`
	RoleID uint `gorm:"primaryKey;not null" json:"role_id"`
}

// Movie model: bảng Movies
type Movie struct {
	gorm.Model
	Title         string `gorm:"size:100;not null;unique" json:"title"`
	Description   string `gorm:"size:255;not null" json:"description"`
	Duration      int    `gorm:"not null" json:"duration"`
	CoverURL      string `gorm:"type:text" json:"cover_url"`
	BackgroundURL string `gorm:"type:text" json:"background_url"`

	// Quan hệ many2many với MovieCategory qua bảng movie_category_mappings
	Categories []MovieCategory `gorm:"many2many:movie_category_mappings;joinForeignKey:MovieID;joinReferences:CategoryID" json:"categories"`
}

// Theater model: bảng Theaters
type Theater struct {
	gorm.Model
	Name string `gorm:"size:10;not null;unique" json:"name"`
}

// Seat model: bảng Seats
type Seat struct {
	gorm.Model
	Type        string  `gorm:"type:text;not null;unique" json:"type"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Number      int     `gorm:"not null" json:"number"`
	EventID     uint    `gorm:"not null" json:"event_id"` // Tham chiếu đến Theater.ID

	// Quan hệ: mỗi Seat thuộc một Theater
	Theater Theater `gorm:"foreignKey:EventID;references:ID" json:"theater"`
}

// Payment model: bảng Payment
type Payment struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id"`
	TotalAmount float64 `gorm:"not null" json:"total_amount"`
	Status      string  `gorm:"size:20;not null;default:pending" json:"status"`

	// Quan hệ: Payment thuộc về User và có nhiều PaymentDetail
	User           User            `gorm:"foreignKey:UserID;references:ID" json:"user"`
	PaymentDetails []PaymentDetail `gorm:"foreignKey:PaymentID" json:"payment_details"`
}

// PaymentDetail model: bảng PaymentDetails
type PaymentDetail struct {
	gorm.Model
	PaymentID  uint    `gorm:"not null" json:"payment_id"`
	SeatID     uint    `gorm:"not null" json:"seat_id"`
	Quantity   int     `gorm:"not null;default:1" json:"quantity"`
	TotalPrice float64 `gorm:"not null" json:"total_price"`

	// Quan hệ: PaymentDetail thuộc về Payment và Seat
	Payment Payment `gorm:"foreignKey:PaymentID;references:ID" json:"payment"`
	Seat    Seat    `gorm:"foreignKey:SeatID;references:ID" json:"seat"`
}

// MovieCategory model: bảng MoviesCategories
type MovieCategory struct {
	gorm.Model
	Name string `gorm:"size:100;not null;unique" json:"name"`
}

// MovieCategoryMapping join table: bảng MovieCategoryMappings
type MovieCategoryMapping struct {
	MovieID    uint `gorm:"primaryKey;not null" json:"movie_id"`
	CategoryID uint `gorm:"primaryKey;not null" json:"category_id"`
}
