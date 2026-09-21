package models

import (
	"time"
	"gorm.io/gorm"
)

type Admin struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      string         `gorm:"type:char(36);not null"`
	User        User
	Name        string         `gorm:"type:varchar(100);not null"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Teacher struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      string         `gorm:"type:char(36);not null"`
	User        User
	NIP         string         `gorm:"type:varchar(50);unique"` // Nomor Induk Pegawai
	Name        string         `gorm:"type:varchar(100);not null"`
	Gender      string         `gorm:"type:enum('L','P')"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Student struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    string         `gorm:"type:char(36);not null"`
	User      User
	NISN      string         `gorm:"type:varchar(50);unique;not null"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Gender    string         `gorm:"type:enum('L','P')"`
	ClassID   uint           `gorm:"not null"`
	Class     *Class         // Pointer to avoid circular init issues before Class is defined
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Curriculum struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      string         `gorm:"type:char(36);not null"`
	User        User
	NIP         string         `gorm:"type:varchar(50);unique"`
	Name        string         `gorm:"type:varchar(100);not null"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Principal struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      string         `gorm:"type:char(36);not null"`
	User        User
	NIP         string         `gorm:"type:varchar(50);unique"`
	Name        string         `gorm:"type:varchar(100);not null"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
