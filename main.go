package main

import (
	"TakeHomeTest/handler"
	"TakeHomeTest/propertyCalculator"
	"log"

	"fmt"
)

func main() {
	// Membuat objek FileProcessor dengan file bernama "file.txt"
	fileProcessor := propertyCalculator.NewFileProcessor("file.txt")

	fileProcessor.WordByWordScan()

	result, err := handler.HandleFileProcessing("file.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Menampilkan hasil yang sudah diurutkan

	log.Println("Data Berdasarkan Luas Tanah Terluas Hingga Terkecilgit ")
	for i, row := range result {
		fmt.Printf("%d: %v\n", i, row)
	}

}
