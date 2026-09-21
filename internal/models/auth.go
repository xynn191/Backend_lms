package models

import (
	"time"
	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"type:varchar(50);not null;unique"` // e.g., admin, guru, siswa, kurikulum, kepala_sekolah
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Users     []User
}

type User struct {
	ID        string         `gorm:"type:char(36);primaryKey"` // UUID
	RoleID    uint           `gorm:"not null"`
	Role      Role
	Email     string         `gorm:"type:varchar(100);not null;unique"`
	Password  string         `gorm:"type:varchar(255);not null"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
