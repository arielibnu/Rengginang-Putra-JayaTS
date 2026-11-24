package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func hubungkanDatabase() error {
	_ = godotenv.Load("backend/.env", "../.env", ".env", "../../backend/.env")

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	pengguna := os.Getenv("DB_USER")
	if pengguna == "" {
		pengguna = "root"
	}
	kataSandi := os.Getenv("DB_PASSWORD")
	namaDatabase := os.Getenv("DB_NAME")
	if namaDatabase == "" {
		namaDatabase = "rengginang_putra_jaya"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=2s",
		pengguna, kataSandi, host, port, namaDatabase,
	)

	koneksi, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi: %w", err)
	}

	koneksi.SetMaxOpenConns(10)
	koneksi.SetMaxIdleConns(5)
	koneksi.SetConnMaxLifetime(5 * time.Minute)

	if err := koneksi.Ping(); err != nil {
		koneksi.Close()
		return fmt.Errorf("gagal terhubung ke MySQL: %w", err)
	}

	db = koneksi
	fmt.Println("Berhasil terhubung ke MySQL!")
	return nil
}
