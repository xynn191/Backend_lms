package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedRoles(db)
	seedAdminUser(db)
	seedActiveAcademicYear(db)
}

func seedActiveAcademicYear(db *gorm.DB) {
	const defaultAcademicYear = "2026/2027 Ganjil"

	var active AcademicYear
	result := db.Where("is_active = ?", true).First(&active)
	if result.Error == nil {
		backfillClassAcademicYears(db, active.ID)
		return
	}

	academicYear := AcademicYear{Name: defaultAcademicYear, IsActive: true}
	if err := db.Where(AcademicYear{Name: defaultAcademicYear}).FirstOrCreate(&academicYear).Error; err != nil {
		log.Printf("[Seed] Gagal membuat tahun ajaran aktif: %v", err)
		return
	}
	if err := db.Model(&AcademicYear{}).Where("id = ?", academicYear.ID).Update("is_active", true).Error; err != nil {
		log.Printf("[Seed] Gagal mengaktifkan tahun ajaran: %v", err)
		return
	}
	backfillClassAcademicYears(db, academicYear.ID)
}

func backfillClassAcademicYears(db *gorm.DB, academicYearID uint) {
	if err := db.Model(&Class{}).Where("academic_year_id = ?", 0).Update("academic_year_id", academicYearID).Error; err != nil {
		log.Printf("[Seed] Gagal memperbarui tahun ajaran kelas lama: %v", err)
	}
}

func seedRoles(db *gorm.DB) {
	roles := []Role{
		{ID: 1, Name: RoleAdmin},
		{ID: 2, Name: RoleTeacher},
		{ID: 3, Name: RoleStudent},
		{ID: 4, Name: RoleCurriculum},
		{ID: 5, Name: RolePrincipal},
	}

	for _, role := range roles {
		var existing Role
		result := db.Where("name = ?", role.Name).First(&existing)
		if result.Error != nil {
			db.Create(&role)
		}
	}
	log.Println("[Seed] Roles seeded.")
}

func seedAdminUser(db *gorm.DB) {
	var existingUser User
	result := db.Where("email = ?", "admin@gmail.com").First(&existingUser)
	if result.Error == nil {
		return // Already exists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("[Seed] Gagal hash password admin:", err)
	}

	adminUser := User{
		ID:       "00000000-0000-0000-0000-000000000001",
		RoleID:   1, // RoleAdmin is seeded with ID 1 for the initial admin.
		Email:    "admin@gmail.com",
		Password: string(hashed),
		IsActive: true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatal("[Seed] Gagal membuat user admin:", err)
	}

	admin := Admin{
		UserID: adminUser.ID,
		Name:   "Super Admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("[Seed] Gagal membuat profil admin:", err)
	}

	log.Println("[Seed] Admin default berhasil dibuat: admin@edulms.id / Admin@123")
}
