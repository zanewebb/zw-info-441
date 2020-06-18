package main

import (
	"os"
	"github.com/exercise-zanewebbUW/zip/zipserver/models"
	"github.com/exercise-zanewebbUW/zip/zipserver/handlers"
)

/*
fileReader = bufio.NewReader(f)
	csvReader = csv.NewReader(fileReader)
	ZipIndex := map[string]int{}
	return ZipSlice
	*/
func main() {
	zipsFile, err := os.Open("zips.csv")
	zSlice, err := models.LoadZips(zipsFile,42613)
	zipsFile.Close()
	zipMap := models.ZipIndex{}
	
	for i := 0; i<len(zSlice); i++ {
			
	}

}