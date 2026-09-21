package models

import (
	"time"
	"gorm.io/gorm"
)

type Major struct {
	ID        uint           `gorm:"primaryKey"`
	Code      string         `gorm:"type:varchar(20);not null;unique"` // e.g., RPL, TKJ
	Name      string         `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Class struct {
	ID                uint           `gorm:"primaryKey"`
	MajorID           uint           `gorm:"not null"`
	Major             Major
	Name              string         `gorm:"type:varchar(50);not null"` // e.g., X-RPL 1
	GradeLevel        string         `gorm:"type:enum('X','XI','XII');not null"`
	HomeroomTeacherID uint           // Wali Kelas
	HomeroomTeacher   *Teacher       `gorm:"foreignKey:HomeroomTeacherID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

type Subject struct {
	ID        uint           `gorm:"primaryKey"`
	Code      string         `gorm:"type:varchar(20);not null;unique"`
	Name      string         `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ClassSubject is the junction table connecting a Class, a Subject, and a Teacher.
type ClassSubject struct {
	ID        uint           `gorm:"primaryKey"`
	ClassID   uint           `gorm:"not null"`
	Class     Class
	SubjectID uint           `gorm:"not null"`
	Subject   Subject
	TeacherID uint           `gorm:"not null"` // The teacher teaching this subject in this class
	Teacher   Teacher
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
