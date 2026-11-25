package models

import (
	"time"

	"gorm.io/gorm"
)

// Role names are persisted in the database and included in JWT claims. Keep
// them in one place so authentication, authorization, and user creation use
// exactly the same values.
const (
	RoleAdmin      = "admin"
	RoleTeacher    = "guru"
	RoleStudent    = "siswa"
	RoleCurriculum = "kurikulum"
	RolePrincipal  = "kepala_sekolah"
	RoleSuperadmin = "superadmin"
)

var SupportedRoles = []string{
	RoleAdmin,
	RoleTeacher,
	RoleStudent,
	RoleCurriculum,
	RolePrincipal,
	RoleSuperadmin,
}

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(50);not null;unique"` // e.g., admin, guru, siswa, kurikulum, kepala_sekolah
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Users     []User
}

type User struct {
	ID        string `gorm:"type:char(36);primaryKey"` // UUID
	RoleID    uint   `gorm:"not null"`
	Role      Role
	Email     string `gorm:"type:varchar(100);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
