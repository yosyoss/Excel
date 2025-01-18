package save

import (
	"Yosyos/src/api/getlist"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"Yosyos/src/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ProcessData(c *gin.Context) {
	excelResponse, err := getlist.ExcelGet(c)
	if err != nil {
		log.Println("Gagal mendapatkan data dari Excel:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data dari Excel"})
		return
	}

	var data []model.TaxData
	for _, item := range excelResponse {
		finalTaxAmt, ok := item["Final Total Amt"].(float64)
		if !ok {
			log.Printf("Failed to assert Final Tax Amt: %v", item["Final Total Amt"])
			finalTaxAmt = 0
		}

		data = append(data, model.TaxData{
			SlipType:    item["Slip Type"].(string),
			StmtNo:      item["Stmt.No."].(string),
			FinalTaxAmt: finalTaxAmt,
		})
	}

	groupedData := make(map[string]float64)
	for _, entry := range data {
		key := fmt.Sprintf("%s-%s", entry.SlipType, entry.StmtNo)
		groupedData[key] += entry.FinalTaxAmt
	}

	var result []model.Response
	for key, finalTaxAmt := range groupedData {
		convertString := strconv.FormatFloat(finalTaxAmt, 'f', 2, 64)
		splitKey := splitKey(key, "-")
		result = append(result, model.Response{
			SlipType:    splitKey[0],
			StmtNo:      splitKey[1],
			FinalTaxAmt: convertString,
		})
	}
	err = saveToExcel(result)
	if err != nil {
		log.Println("Gagal menyimpan ke Excel:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file Excel"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func splitKey(key string, delimiter string) []string {
	return split(key, delimiter)
}

func split(input string, delimiter string) []string {
	splitResult := make([]string, 0)
	temp := ""
	for _, char := range input {
		if string(char) == delimiter {
			splitResult = append(splitResult, temp)
			temp = ""
		} else {
			temp += string(char)
		}
	}
	splitResult = append(splitResult, temp)
	return splitResult
}

func saveToExcel(data []model.Response) error {
	folderPath := "src/DataExcel/Result"
	fileName := "processed_data.xlsx"
	filePath := fmt.Sprintf("%s/%s", folderPath, fileName)

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("gagal membuat folder: %w", err)
	}

	f := excelize.NewFile()

	f.NewSheet("Sheet1")

	headers := []string{"Slip Type", "Stmt.No.", "Final Tax Amt"}
	for col, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+col)))
		f.SetCellValue("Sheet1", cell, header)
	}

	for i, item := range data {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), item.SlipType)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), item.StmtNo)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), item.FinalTaxAmt)
	}

	err = f.SaveAs(filePath)
	if err != nil {
		return fmt.Errorf("gagal menyimpan file Excel: %w", err)
	}

	log.Printf("File berhasil disimpan ke: %s", filePath)
	return nil
}
