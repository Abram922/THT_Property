package handler


import "TakeHomeTest/propertyCalculator"



func HandleFileProcessing(fileName string) ([][]string, error) {
    // Memanggil fungsi di service untuk membaca dan mengurutkan data
    return propertyCalculator.ProcessAndSortFile(fileName)
}