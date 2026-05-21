package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

type SVRResult struct {
	PrediksiSkorQ5 float64 `json:"prediksi_skor_q5"`
	Status         string  `json:"status"`
	DataPointsUsed int     `json:"data_points_used"`
	MarginOfError  float64 `json:"margin_of_error"`
	Error          string  `json:"error"`
}

var dataAssesmen []Assesment
var idCounter int = 1

var daftarPertanyaan = [5]string{
	"Kurangnya minat atau kesenangan dalam melakukan aktivitas sehari-hari",
	"Merasa down, sedih, tertekan, atau putus asa",
	"Merasa gugup, cemas, gelisah, atau seolah-olah sesuatu yang buruk akan terjadi",
	"Merasa kesulitan untuk menghentikan atau mengendalikan rasa khawatir",
	"Seberapa besar masalah-masalah di atas mengganggu aktivitas Anda (seperti kuliah/kerja, urusan rumah, atau bersosialisasi)?",
}

func main() {
	loadFile()
	var pilihan int

	for {
		cetakHeader("SISTEM SELF-ASSESSMENT")
		fmt.Println("1. Tambah Hasil Assessment")
		fmt.Println("2. Ubah Hasil Assessment")
		fmt.Println("3. Hapus Hasil Assessment")
		fmt.Println("4. Cari Assessment Berdasarkan User ID")
		fmt.Println("5. Urutkan Data Assessment")
		fmt.Println("6. Tampilkan Laporan Pengguna")
		fmt.Println("0. Keluar")
		cetakGaris("-", 60)
		fmt.Print("Pilih menu: ")
		fmt.Scan(&pilihan)

		switch pilihan {
		case 1:
			tambahData()
		case 2:
			ubahData()
		case 3:
			hapusData()
		case 4:
			menuCari()
		case 5:
			menuUrutkan()
		case 6:
			tampilkanLaporan()
		case 0:
			fmt.Println("Terima kasih telah menggunakan sistem ini.")
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}
	}
}

func cetakGaris(karakter string, n int) {
	fmt.Println(strings.Repeat(karakter, n))
}

func cetakHeader(judul string) {
	const jumKar int = 60
	fmt.Println()
	cetakGaris("=", jumKar)
	spasi := (jumKar - len(judul)) / 2
	fmt.Print(strings.Repeat(" ", spasi))
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

func tentukanRekomendasi(skor int) string {
	switch {
	case skor >= 22:
		return "Pertahankan pola hidup dan kebiasaan positif Anda."
	case skor >= 17:
		return "Kondisi baik, tetap jaga keseimbangan antara aktivitas dan istirahat."
	case skor >= 12:
		return "Perbanyak relaksasi dan perbaiki pola tidur Anda."
	case skor >= 7:
		return "Kurangi beban kerja, pertimbangkan untuk bercerita pada orang terdekat."
	default:
		return "Segera konsultasikan kondisi Anda kepada tenaga profesional (psikolog/dokter)."
	}
}

func tambahData() {
	var a Assesment
	a.ID = idCounter

	fmt.Print("Masukkan User ID: ")
	fmt.Scan(&a.UserID)

	fmt.Print("Masukkan Tanggal (YYYY-MM-DD): ")
	fmt.Scan(&a.Tanggal)

	fmt.Println("\nDalam 2 minggu terakhir, seberapa sering Anda terganggu oleh masalah-masalah berikut?")
	fmt.Println("(Berikan nilai skala 1-5)")
	cetakGaris("-", 60)

	for i := 0; i < 5; i++ {
		fmt.Printf("%d. %s\n", i+1, daftarPertanyaan[i])
		fmt.Print("   Jawaban (1-5): ")
		fmt.Scan(&a.Jawaban[i])
	}

	a.SkorTotal = hitungSkor(a.Jawaban)
	a.Kategori = tentukanKategori(a.SkorTotal)
	a.Rekomendasi = tentukanRekomendasi(a.SkorTotal)

	dataAssesmen = append(dataAssesmen, a)
	idCounter++

	saveFile()
	fmt.Println("\nData berhasil ditambahkan!")
}

func ubahData() {
	var id, idxUbah int
	ditemukan := false

	fmt.Print("Masukkan ID Assessment yang ingin diubah: ")
	fmt.Scan(&id)

	for i, v := range dataAssesmen {
		if v.ID == id {
			idxUbah = i
			ditemukan = true
			break
		}
	}

	if !ditemukan {
		fmt.Println("Data tidak ditemukan.")
		return
	}

	fmt.Print("Masukkan Tanggal Baru (YYYY-MM-DD): ")
	fmt.Scan(&dataAssesmen[idxUbah].Tanggal)

	fmt.Println("\nDalam 2 minggu terakhir, seberapa sering Anda terganggu oleh masalah-masalah berikut?")
	fmt.Println("(Berikan nilai skala 1-5 baru)")
	cetakGaris("-", 60)

	for i := 0; i < 5; i++ {
		fmt.Printf("%d. %s\n", i+1, daftarPertanyaan[i])
		fmt.Print("   Jawaban Baru (1-5): ")
		fmt.Scan(&dataAssesmen[idxUbah].Jawaban[i])
	}

	dataAssesmen[idxUbah].SkorTotal = hitungSkor(dataAssesmen[idxUbah].Jawaban)
	dataAssesmen[idxUbah].Kategori = tentukanKategori(dataAssesmen[idxUbah].SkorTotal)
	dataAssesmen[idxUbah].Rekomendasi = tentukanRekomendasi(dataAssesmen[idxUbah].SkorTotal)

	saveFile()

	fmt.Println("\nData berhasil diperbarui!")
}

func hapusData() {
	var id int
	fmt.Print("Masukkan ID Assessment yang ingin dihapus: ")
	fmt.Scan(&id)

	for i, v := range dataAssesmen {
		if v.ID == id {
			dataAssesmen = append(dataAssesmen[:i], dataAssesmen[i+1:]...)

			saveFile()

			fmt.Println("Data berhasil dihapus!")
			return
		}
	}
	fmt.Println("Data tidak ditemukan.")
}

func menuCari() {
	var pilihan, userID int
	fmt.Println("1. Sequential Search")
	fmt.Println("2. Binary Search")
	fmt.Print("Pilih metode pencarian: ")
	fmt.Scan(&pilihan)

	fmt.Print("Masukkan User ID yang dicari: ")
	fmt.Scan(&userID)

	var hasil []Assesment
	if pilihan == 1 {
		hasil = sequentialSearch(userID)
	} else if pilihan == 2 {
		hasil = binarySearch(userID)
	}

	if len(hasil) > 0 {
		cetakHeader("Hasil Pencarian")
		for _, v := range hasil {
			fmt.Printf("ID: %d | Skor: %d | Tgl: %s | Kategori: %s\n", v.ID, v.SkorTotal, v.Tanggal, v.Kategori)
		}
	} else {
		fmt.Println("Data assessment untuk User ID tersebut tidak ditemukan.")
	}
}

func sequentialSearch(userID int) []Assesment {
	var hasil []Assesment
	for _, v := range dataAssesmen {
		if v.UserID == userID {
			hasil = append(hasil, v)
		}
	}
	return hasil
}

func binarySearch(userID int) []Assesment {
	tempData := make([]Assesment, len(dataAssesmen))
	copy(tempData, dataAssesmen)

	for i := 0; i < len(tempData)-1; i++ {
		for j := 0; j < len(tempData)-i-1; j++ {
			if tempData[j].UserID > tempData[j+1].UserID {
				tempData[j], tempData[j+1] = tempData[j+1], tempData[j]
			}
		}
	}

	var hasil []Assesment
	low, high := 0, len(tempData)-1

	ditemukanIdx := -1
	for low <= high {
		mid := low + (high-low)/2
		if tempData[mid].UserID == userID {
			ditemukanIdx = mid
			break
		} else if tempData[mid].UserID < userID {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	if ditemukanIdx != -1 {
		hasil = append(hasil, tempData[ditemukanIdx])
		for i := ditemukanIdx - 1; i >= 0 && tempData[i].UserID == userID; i-- {
			hasil = append(hasil, tempData[i])
		}
		for i := ditemukanIdx + 1; i < len(tempData) && tempData[i].UserID == userID; i++ {
			hasil = append(hasil, tempData[i])
		}
	}

	return hasil
}

func menuUrutkan() {
	var pilihan int
	fmt.Println("1. Urutkan berdasarkan Skor Total (Selection Sort - Descending)")
	fmt.Println("2. Urutkan berdasarkan Tanggal (Insertion Sort - Terbaru ke Terlama)")
	fmt.Print("Pilih pengurutan: ")
	fmt.Scan(&pilihan)

	if pilihan == 1 {
		selectionSortSkor()
		fmt.Println("Data berhasil diurutkan berdasarkan Skor Tertinggi.")
	} else if pilihan == 2 {
		insertionSortTanggal()
		fmt.Println("Data berhasil diurutkan berdasarkan Tanggal Terbaru.")
	}
}

func selectionSortSkor() {
	n := len(dataAssesmen)
	for i := 0; i < n-1; i++ {
		maxIdx := i
		for j := i + 1; j < n; j++ {
			if dataAssesmen[j].SkorTotal > dataAssesmen[maxIdx].SkorTotal {
				maxIdx = j
			}
		}
		dataAssesmen[i], dataAssesmen[maxIdx] = dataAssesmen[maxIdx], dataAssesmen[i]
	}
}

func insertionSortTanggal() {
	n := len(dataAssesmen)
	for i := 1; i < n; i++ {
		key := dataAssesmen[i]
		j := i - 1
		for j >= 0 && dataAssesmen[j].Tanggal < key.Tanggal {
			dataAssesmen[j+1] = dataAssesmen[j]
			j--
		}
		dataAssesmen[j+1] = key
	}
}

func tampilkanLaporan() {
	var userID int
	fmt.Print("Masukkan User ID untuk mencetak laporan: ")
	fmt.Scan(&userID)

	historiSemua := sequentialSearch(userID)

	if len(historiSemua) == 0 {
		fmt.Println("Tidak ada data untuk User ID ini.")
		return
	}

	var historiUser []Assesment
	totalSkor := 0
	batasWaktu := time.Now().AddDate(0, 0, -30)

	for _, v := range historiSemua {
		tglData, err := time.Parse("2006-01-02", v.Tanggal)
		if err != nil {
			historiUser = append(historiUser, v)
			continue
		}

		if tglData.After(batasWaktu) || tglData.Equal(batasWaktu) {
			historiUser = append(historiUser, v)
		}
	}

	if len(historiUser) == 0 {
		fmt.Println("\n[Notifikasi] User ini memiliki riwayat lama, tetapi tidak ada aktivitas assessment dalam 1 bulan terakhir.")
		return
	}

	cetakHeader(fmt.Sprintf("LAPORAN AKTIVITAS 1 BULAN TERAKHIR (USER ID: %d)", userID))


	batas := len(historiUser)
	if batas > 5 {
		batas = 5
	}

	fmt.Printf("Menampilkan %d riwayat terakhir pada bulan ini:\n", batas)
	for i := len(historiUser) - 1; i >= len(historiUser)-batas; i-- {
		v := historiUser[i]
		fmt.Printf("- Tgl: %s | Skor: %d | Kategori: %s\n", v.Tanggal, v.SkorTotal, v.Kategori)
		fmt.Printf("  Rekomendasi: %s\n", v.Rekomendasi)
	}

	for _, v := range historiUser {
		totalSkor += v.SkorTotal
	}
	rataRata := float64(totalSkor) / float64(len(historiUser))

	cetakGaris("-", 60)
	fmt.Printf("Total Assessment (1 Bulan Terakhir): %d kali\n", len(historiUser))
	fmt.Printf("Rata-rata Skor Sebulan Terakhir    : %.2f\n", rataRata)
	cetakGaris("-", 60)


	asesmenTerbaru := historiUser[len(historiUser)-1]
	empatJawaban := asesmenTerbaru.Jawaban[0:4]

	jawabanJSON, err := json.Marshal(empatJawaban)
	if err != nil {
		return
	}

	fmt.Println("\n[AI System] Menganalisis kecenderungan dampak produktivitas sebulan terakhir...")

	cmd := exec.Command("python", "svr_worker.py", string(jawabanJSON))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return
	}

	var hasil SVRResult
	err = json.Unmarshal(out.Bytes(), &hasil)
	if err != nil {
		return
	}

	if hasil.Error == "" {
		fmt.Println("\n=============================================")
		fmt.Println("     ANALISIS PREDIKSI AI (MODEL SVR)       ")
		fmt.Println("=============================================")
		fmt.Printf("Model SVR dilatih dari %d data global sebulan terakhir.\n", hasil.DataPointsUsed)
		fmt.Printf("Prediksi Dampak Aktivitas User     : %.2f \u00b1 %.2f (Skala 1-5)\n", hasil.PrediksiSkorQ5, hasil.MarginOfError)
		fmt.Printf("Jawaban Aktual (Realitas) User     : %d\n", asesmenTerbaru.Jawaban[4])
		
		selisih := hasil.PrediksiSkorQ5 - float64(asesmenTerbaru.Jawaban[4])
		if selisih < 0 {
			selisih = -selisih
		}
		
		if selisih <= 0.5 {
			fmt.Println("Insight: Tingkat produktivitas user ini sangat konsisten dengan tren data sebulan terakhir.")
		} else if hasil.PrediksiSkorQ5 > float64(asesmenTerbaru.Jawaban[4]) {
			fmt.Println("Insight: Gejala klinis user berpotensi menekan produktivitas lebih tinggi dibanding yang disadari.")
		} else {
			fmt.Println("Insight: User menunjukkan resiliensi atau mekanisme koping aktivitas yang sangat baik.")
		}
		fmt.Println("=============================================")
	}
}

func saveFile() {
	file, err := os.Create("data_assesmen.json")
	if err != nil {
		fmt.Println("Gagal membuat file penyimpanan:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dataAssesmen)
	if err != nil {
		fmt.Println("Gagal menyimpan data:", err)
	}
}

func loadFile() {
	file, err := os.Open("data_assesmen.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Gagal membuka file penyimpanan:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dataAssesmen)
	if err != nil {
		fmt.Println("Gagal membaca data:", err)
		return
	}

	for _, v := range dataAssesmen {
		if v.ID >= idCounter {
			idCounter = v.ID + 1
		}
	}
}