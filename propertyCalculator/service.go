package propertyCalculator

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var priceAddsCornerLot float32
var priceBeforeBonus float32
var PricelandAppreciates float32
var PriceLandDepreciates float32
var propertyPrice float32
var propertyFee float32
var standardLoc int = 10000000
var premiumLoc int = 15000000
var landAppreciates float32 = 5
var residentalBuildingDepreciates float32 = 2.5
var commercialBuildingDepreciates float32 = 3.5
var bonusPremiumLoc int = 20
var bonusCornerLoc int = 15
var feeResidentalBase int = 2500
var feeCommercialbas int = 3500
var feeSecurity int = 1000
var feeCleaning int = 800

type FileProcessor struct {
	FilePath string
	Property *Property
}

func NewFileProcessor(filename string) *FileProcessor {
	// Mendapatkan direktori saat ini
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Membuat path lengkap file
	filePath := filepath.Join(currentDir, filename)

	// Membuat objek FileProcessor dengan properti yang belum di-set
	return &FileProcessor{
		FilePath: filePath,
		Property: &Property{}, // Inisialisasi entitas Property
	}
}

func (fp *FileProcessor) WordByWordScan() {
	file, err := os.Open(fp.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	infoFile, err := file.Stat()
	if err != nil {
		log.Fatalf("Gagal mendapatkan informasi file: %v", err)
	}

	if infoFile.Size() == 0 {
		log.Println("Blank File")
	}

	scanner := bufio.NewScanner(file)

	count := 0
	propertyValues := []interface{}{} // Slice untuk menampung nilai yang akan di-set

	// Loop untuk membaca baris per baris dari file
	for scanner.Scan() {
		line := scanner.Text() // Membaca satu baris
		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords) // Memindai kata dalam baris

		// Loop untuk membaca kata per kata dalam satu baris
		for wordScanner.Scan() {
			word := wordScanner.Text()
			//fmt.Println(word)

			// Menambahkan setiap kata yang dipindai ke dalam slice
			propertyValues = append(propertyValues, word)

			// Jika sudah mengumpulkan 8 nilai, kita akan set nilai ke Property
			if len(propertyValues) == 8 {
				// Set nilai ke Property sesuai urutan
				fp.Property.SetPropertyValues(propertyValues)
				propertyValues = []interface{}{} // Reset untuk iterasi berikutnya
			}

			count++

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

// Fungsi untuk mengisi nilai ke dalam Property
func (p *Property) SetPropertyValues(values []interface{}) {

	re := regexp.MustCompile(`\d+`)
	// Pastikan jumlah nilai sesuai dengan jumlah properti

	if len(values) != 8 {
		log.Fatal("Jumlah nilai tidak sesuai dengan jumlah properti")
		return
	}

	// Menetapkan nilai ke properti sesuai urutan
	var err error

	// Timestamp
	p.Timestamp, err = time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%v", values[0]))
	if err != nil {
		log.Printf("Error parsing Timestamp: %v\n", err)
	}
	// Type
	p.Type = fmt.Sprintf("%v", values[1])
	// Area
	p.Area, err = parseFloat(values[2])
	if err != nil {
		log.Printf("Error parsing Area: %v\n", err)
	} else if p.Area <= 0 {
		log.Printf("Error Area Tanah Lebih Kecil atau Sama Dengan 0\n")
	}
	// BuildYear
	YearStr := re.FindString(fmt.Sprintf("%v", values[3])) // Ambil angka dari nilai di values[6]
	p.BuildYear, err = strconv.Atoi(YearStr)               // Parsing nilai angka
	if err != nil {
		log.Printf("Error parsing Build Year: %v\n", err)
	} else if p.BuildYear > time.Now().Year() {
		log.Printf("Kesalahan Tahun, Tahun berada di masa depan  \n")
	} else if p.BuildYear < 1900 {
		log.Printf("Kesalahan Tahun, Tahun tidak valid  \n")
	}
	// Location
	p.Location = fmt.Sprintf("%v", values[4])
	if p.Location != "PREMIUM" && p.Location != "STANDARD" {
		log.Println("Lokasi Tidak Valid")
	}
	// Corner
	p.Corner = fmt.Sprintf("%v", values[5])
	if p.Corner != "YES" && p.Corner != "NO" && p.Corner != "CORNER" {
		log.Println("Corner Tidak Valid")
	}
	// Parking
	parkingStr := re.FindString(fmt.Sprintf("%v", values[6])) // Ambil angka dari nilai di values[6]
	p.Parking, err = strconv.Atoi(parkingStr)                 // Parsing nilai angka
	if err != nil {
		log.Printf("Error parsing Parking: %v\n", err)
	} else if p.Parking < 0 || p.Parking > 99 {
		log.Println("Parking Tidak Valid")

	}
	// Facilities
	p.Facilities = []string{fmt.Sprintf("%v\n\n", values[7])}

	if p.Location == "STANDARD" {

		pangkat := time.Now().Year() - p.BuildYear
		if pangkat == 0 {
			pangkat = pangkat + 1
		} else {
			pangkat = pangkat + 0
		}

		PricelandAppreciates := float32(standardLoc) * p.Area * float32(math.Pow(1+float64(landAppreciates)/100, float64(pangkat)))

		if p.Type == "COMMERCIAL" {
			PriceLandDepreciates = float32(standardLoc) * p.Area * float32(math.Pow(1+float64(residentalBuildingDepreciates)/100, float64(pangkat)))
			propertyPrice = float32(PricelandAppreciates + PriceLandDepreciates)
		} else if p.Type == "RESIDENTIAL" {
			PriceLandDepreciates = float32(standardLoc) * p.Area * float32(math.Pow(1+float64(commercialBuildingDepreciates)/100, float64(pangkat)))
			propertyPrice = PricelandAppreciates + PriceLandDepreciates
		} else {
			propertyPrice = propertyPrice + 0
		}

		if p.Corner == "YES" || p.Corner == "CORNER" {
			priceAddsCornerLot := propertyPrice * (float32(bonusCornerLoc) / 100)
			propertyPrice = propertyPrice + priceAddsCornerLot
		} else {
			propertyPrice = propertyPrice + 0
		}

		propertyFee = float32(feeCleaning+feeCommercialbas+feeResidentalBase+feeSecurity) * (p.Area)

	} else if p.Location == "PREMIUM" {
		pangkat := time.Now().Year() - p.BuildYear
		if pangkat == 0 {
			pangkat = pangkat + 1
		} else {
			pangkat = pangkat + 0
		}

		PricelandAppreciates = float32(premiumLoc) * p.Area * float32(math.Pow(1+float64(landAppreciates)/100, float64(pangkat)))

		if p.Type == "COMMERCIAL" {
			PriceLandDepreciates = float32(premiumLoc) * p.Area * float32(math.Pow(1+float64(residentalBuildingDepreciates)/100, float64(pangkat)))
			priceBeforeBonus = PricelandAppreciates + PriceLandDepreciates

		} else if p.Type == "RESIDENTIAL" {
			PriceLandDepreciates = float32(premiumLoc) * p.Area * float32(math.Pow(1+float64(commercialBuildingDepreciates)/100, float64(pangkat)))
			priceBeforeBonus = PricelandAppreciates + PriceLandDepreciates
		} else {
			priceBeforeBonus = PricelandAppreciates + 0
		}

		premiumadds := priceBeforeBonus * float32(bonusPremiumLoc) / 100

		if p.Corner == "YES" || p.Corner == "CORNER" {
			priceAddsCornerLot = propertyPrice * (float32(bonusCornerLoc) / 100)
		} else {
			priceAddsCornerLot = 0
		}

		propertyPrice = premiumadds + priceAddsCornerLot

		propertyFee = float32(feeCleaning+feeCommercialbas+feeResidentalBase+feeSecurity) * (p.Area)

		// tambahkan bonus lokasi premium

	}

	fmt.Printf("Property Value: Rp %s\n", humanize.Comma(int64(propertyPrice)))
	fmt.Printf("Monthly Maintenance: Rp %s\n\n", humanize.Comma(int64(propertyFee)))

}

// Helper function untuk parsing float
func parseFloat(value interface{}) (float32, error) {
	switch v := value.(type) {
	case float32:
		return v, nil
	case string:
		var f float32
		// Memperbaiki penggunaan fmt.Sscanf untuk menangkap nilai float
		_, err := fmt.Sscanf(v, "%f", &f)
		if err != nil {
			return 0, fmt.Errorf("Error parsing float: %v", err)
		}
		return f, nil
	}
	return 0, fmt.Errorf("Invalid float value: %v", value)
}

// Helper function untuk parsing int

func ProcessAndSortFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Fields(line)
		data = append(data, columns)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(data, func(i, j int) bool {
		val1, _ := strconv.ParseFloat(data[i][2], 64)
		val2, _ := strconv.ParseFloat(data[j][2], 64)
		return val1 > val2 // Dari terbesar ke terkecil

	})

	return data, nil

}
