package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("File .env tidak ditemukan")
	}

	konfigurasi := mysql.NewConfig()
	konfigurasi.User = os.Getenv("DB_USER")
	konfigurasi.Passwd = os.Getenv("DB_PASSWORD")
	konfigurasi.Net = "tcp"
	konfigurasi.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	konfigurasi.DBName = os.Getenv("DB_NAME")
	konfigurasi.ParseTime = true

	db, err := sql.Open("mysql", konfigurasi.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Koneksi database gagal: ", err)
	}

	pembaca := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan nama pengguna: ")
	namaPengguna, _ := pembaca.ReadString('\n')
	namaPengguna = strings.TrimSpace(namaPengguna)

	if namaPengguna == "" {
		log.Fatal("Nama pengguna tidak boleh kosong")
	}

	fmt.Print("Masukkan kata sandi (minimal 8 karakter): ")
	kataSandi, _ := pembaca.ReadString('\n')
	kataSandi = strings.TrimSpace(kataSandi)

	if len(kataSandi) < 8 || len(kataSandi) > 72 {
		log.Fatal("Kata sandi harus 8 sampai 72 byte")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(kataSandi),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal("Gagal membuat hash: ", err)
	}

	_, err = db.Exec(
		`INSERT INTO akun_pemilik
		(nama_pengguna, kata_sandi_hash)
		VALUES (?, ?)`,
		namaPengguna,
		string(hash),
	)
	if err != nil {
		log.Fatal("Gagal menyimpan akun: ", err)
	}

	fmt.Println("Akun pemilik berhasil dibuat!")
}
