package models

import (
	"time"
	"gorm.io/gorm"
)

type Material struct {
	ID             uint           `gorm:"primaryKey"`
	ClassSubjectID uint           `gorm:"not null"`
	ClassSubject   ClassSubject
	Title          string         `gorm:"type:varchar(200);not null"`
	Content        string         `gorm:"type:text"`
	AttachmentURL  string         `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Task struct {
	ID             uint           `gorm:"primaryKey"`
	ClassSubjectID uint           `gorm:"not null"`
	ClassSubject   ClassSubject
	Title          string         `gorm:"type:varchar(200);not null"`
	Description    string         `gorm:"type:text"`
	DueDate        time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TaskSubmission struct {
	ID          uint           `gorm:"primaryKey"`
	TaskID      uint           `gorm:"not null"`
	Task        Task
	StudentID   uint           `gorm:"not null"`
	Student     Student
	FileURL     string         `gorm:"type:varchar(255)"`
	Score       *float64       // Pointer so it can be null if not graded yet
	SubmittedAt time.Time
	GradedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Assessment struct {
	ID              uint           `gorm:"primaryKey"`
	ClassSubjectID  uint           `gorm:"not null"`
	ClassSubject    ClassSubject
	Title           string         `gorm:"type:varchar(200);not null"`
	StartTime       time.Time
	EndTime         time.Time
	DurationMinutes int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type AssessmentSubmission struct {
	ID           uint           `gorm:"primaryKey"`
	AssessmentID uint           `gorm:"not null"`
	Assessment   Assessment
	StudentID    uint           `gorm:"not null"`
	Student      Student
	Score        *float64
	SubmittedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
