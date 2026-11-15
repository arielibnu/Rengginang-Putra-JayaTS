package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Membuat hash kata sandi sebelum disimpan ke database
func buatHashKataSandi(kataSandi string) (string, error) {
	if len(kataSandi) < 8 {
		return "", fmt.Errorf("kata sandi minimal 8 karakter")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(kataSandi),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Memeriksa kata sandi terhadap hash yang tersimpan
func periksaKataSandi(hash, kataSandi string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(kataSandi),
	)

	return err == nil
}
