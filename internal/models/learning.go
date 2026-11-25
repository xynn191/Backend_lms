package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	ID             uint `gorm:"primaryKey"`
	ClassSubjectID uint `gorm:"not null"`
	ClassSubject   ClassSubject
	Title          string `gorm:"type:varchar(200);not null"`
	Content        string `gorm:"type:text"`
	AttachmentURL  string `gorm:"type:text"`
	IsPublished    bool   `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Task struct {
	ID             uint `gorm:"primaryKey"`
	ClassSubjectID uint `gorm:"not null"`
	ClassSubject   ClassSubject
	Title          string    `gorm:"type:varchar(200);not null"`
	Description    string    `gorm:"type:text"`
	AttachmentURL  string    `gorm:"type:text"`
	AttachmentLink string    `gorm:"type:varchar(500)"`
	DueDate        time.Time
	IsPublished    bool      `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TaskSubmission struct {
	ID          uint `gorm:"primaryKey"`
	TaskID      uint `gorm:"not null"`
	Task        Task
	StudentID   uint `gorm:"not null"`
	Student     Student
	FileURL     string     `gorm:"type:varchar(255)"`
	LinkURL     string     `gorm:"type:varchar(500)"`
	Score       *float64   // Pointer so it can be null if not graded yet
	Feedback    string     `gorm:"type:text"`
	SubmittedAt time.Time
	GradedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Assessment struct {
	ID              uint `gorm:"primaryKey"`
	ClassSubjectID  uint `gorm:"not null"`
	ClassSubject    ClassSubject
	Title           string    `gorm:"type:varchar(200);not null"`
	Description     string    `gorm:"type:text"`
	StartTime       time.Time
	EndTime         time.Time
	DurationMinutes int
	IsPublished     bool                 `gorm:"default:false"`
	Questions       []AssessmentQuestion `gorm:"foreignKey:AssessmentID;constraint:OnDelete:CASCADE;"`
	Submissions     []AssessmentSubmission `gorm:"foreignKey:AssessmentID;constraint:OnDelete:CASCADE;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type AssessmentQuestion struct {
	ID           uint `gorm:"primaryKey"`
	AssessmentID uint `gorm:"not null"`
	QuestionText string `gorm:"type:text;not null"`
	Options      []AssessmentOption `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type AssessmentOption struct {
	ID         uint   `gorm:"primaryKey"`
	QuestionID uint   `gorm:"not null"`
	OptionText string `gorm:"type:text;not null"`
	IsCorrect  bool   `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type AssessmentSubmission struct {
	ID           uint `gorm:"primaryKey"`
	AssessmentID uint `gorm:"not null"`
	Assessment   Assessment
	StudentID    uint `gorm:"not null"`
	Student      Student
	Score        *float64
	SubmittedAt  time.Time
	Answers      []AssessmentAnswer `gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type AssessmentAnswer struct {
	ID               uint `gorm:"primaryKey"`
	SubmissionID     uint `gorm:"not null"`
	QuestionID       uint `gorm:"not null"`
	SelectedOptionID uint `gorm:"not null"`
	IsCorrect        bool `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
