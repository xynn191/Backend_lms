package main

import (
	"fmt"
	"log"

	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Gagal load konfigurasi:", err)
	}
	fmt.Println("Memulai server EduLMS Backend...")

	// 2. Koneksi ke database
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}
	fmt.Println("Koneksi ke database berhasil!")

	// 3. Auto-migrate semua tabel
	if err := models.AutoMigrate(db); err != nil {
		log.Fatal("Gagal auto-migrate:", err)
	}
	fmt.Println("Auto-migrate berhasil!")

	// 4. Seed data awal (roles + admin default)
	models.Seed(db)

	// 5. Setup router dan jalankan server
	r := router.SetupRouter(db, cfg)
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("Server berjalan di http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
