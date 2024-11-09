package handler

import "TakeHomeTest/propertyCalculator"

type FileHandler struct {
	processor *propertyCalculator.FileProcessor
}

func NewFileHandler(sfilename string) *FileHandler {
	processor := propertyCalculator.NewFileProcessor(sfilename)
	return &FileHandler{processor: processor}
}

func (fh *FileHandler) ProcessFile() {
	fh.processor.WordByWordScan()
}
