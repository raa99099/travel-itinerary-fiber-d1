package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:120;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"size:20;default:user;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Itinerary struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string    `json:"title" gorm:"size:150;not null"`
	Destination string    `json:"destination" gorm:"size:150;not null"`
	StartDate   string    `json:"start_date" gorm:"size:30;not null"`
	EndDate     string    `json:"end_date" gorm:"size:30;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Spot struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ItineraryID   uint      `json:"itinerary_id" gorm:"not null"`
	Itinerary     Itinerary `json:"itinerary,omitempty" gorm:"foreignKey:ItineraryID"`
	Name          string    `json:"name" gorm:"size:150;not null"`
	Category      string    `json:"category" gorm:"size:80;not null"`
	VisitStatus   string    `json:"visit_status" gorm:"size:30;default:planned;not null"`
	EstimatedCost int64     `json:"estimated_cost" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SpotNote struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	SpotID       uint      `json:"spot_id" gorm:"not null"`
	Spot         Spot      `json:"spot,omitempty" gorm:"foreignKey:SpotID"`
	Title        string    `json:"title" gorm:"size:150;not null"`
	NoteContent  string    `json:"note_content" gorm:"type:text;not null"`
	OpeningHours string    `json:"opening_hours" gorm:"size:80"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
