package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedRoles(db)
	seedAdminUser(db)
}

func seedRoles(db *gorm.DB) {
	roles := []Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "guru"},
		{ID: 3, Name: "siswa"},
		{ID: 4, Name: "kurikulum"},
		{ID: 5, Name: "kepala_sekolah"},
	}

	for _, role := range roles {
		var existing Role
		result := db.First(&existing, role.ID)
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
		RoleID:   1, // admin
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
