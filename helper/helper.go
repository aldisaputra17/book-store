package helper

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomISBN() string {
	sourcer := rand.NewSource(time.Now().UnixNano())
	random := rand.New(sourcer)

	// Menghasilkan 9 angka acak untuk digit ISBN pertama
	isbnDigits := make([]int, 9)
	for i := 0; i < 9; i++ {
		isbnDigits[i] = random.Intn(10)
	}

	// Menghitung digit ke-10 (digit checksum) berdasarkan digit ISBN pertama
	checksum := CalculateISBNChecksum(isbnDigits)

	// Menggabungkan digit ISBN pertama, tanda "-", dan digit checksum menjadi string ISBN
	isbn := JoinISBNDigits(isbnDigits, checksum)

	return isbn
}

func CalculateISBNChecksum(digits []int) int {
	checksum := 0
	for i, digit := range digits {
		checksum += (i + 1) * digit
	}
	checksum %= 11

	return checksum
}

func JoinISBNDigits(digits []int, checksum int) string {
	isbn := "978-"
	for _, digit := range digits {
		isbn += strconv.Itoa(digit)
	}
	isbn += "-" + strconv.Itoa(checksum)

	return isbn
}
