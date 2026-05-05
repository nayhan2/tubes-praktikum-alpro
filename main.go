package main

import (
	"fmt"
)

type Assesment struct {
	ID int
	UserID int
	Jawaban [5]int
	Tanggal string
	SkorTotal int
	Kategori string
	Rekomendasi string
}

// Variable global
var dataAssesmen []Assesment
var idCounter int = 1

func main() {
	fmt.Println("gemas")
}
