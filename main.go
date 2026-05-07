package main

import (
	"fmt"
)

type Assesment struct {
	ID          int
	UserID      int
	Jawaban     [5]int
	Tanggal     string
	SkorTotal   int
	Kategori    string
	Rekomendasi string
}

// Variable global
var dataAssesmen []Assesment
var idCounter int = 1

func main() {
	fmt.Println("gemas")
}

// utility
func cetakGaris(karakter string, n int) {
	for i := 0; i < n; i++ {
		fmt.Print(karakter)
	}
	fmt.Println()
}

func cetakHeader(judul string) {
	const jumKar int = 60
	cetakGaris("=", jumKar)
	n := len(judul)
	spasi := (jumKar - n) / 2
	for i := 0; i < spasi; i++ {
		fmt.Print(" ")
	}
	fmt.Println(judul)
	cetakGaris("=", jumKar)
}

func hitungSkor(ans [5]int) int {
	total := 0
	for i := 0; i < len(ans); i++ {
		total += ans[i]
	}
	return total
}

// program func

func tentukanKategori(skor int) string {
	switch {
	case skor >= 22:
		return "Sangat Baik"
	case skor >= 17:
		return "Baik"
	case skor >= 12:
		return "Cukup"
	case skor >= 7:
		return "Perlu Perhatian"
	default:
		return "Butuh Bantuan Segera"
	}
}
