package getlist

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetListExcelValue(c *gin.Context) {
	resultChannel := make(chan struct {
		data interface{}
		err  error
	})

	go func() {
		data, err := ExcelGet(c)
		resultChannel <- struct {
			data interface{}
			err  error
		}{data, err}
	}()

	result := <-resultChannel

	if result.err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal memproses data Excel",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": result.data,
	})
}

func ExcelGet(c *gin.Context) ([]map[string]interface{}, error) {
	f, err := excelize.OpenFile("src/DataExcel/test.xlsx")
	if err != nil {
		log.Println("Gagal membuka file Excel:", err)
		return nil, fmt.Errorf("gagal membuka file Excel: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Gagal menutup file Excel:", err)
		}
	}()

	rows, err := f.Rows("Sheet1")
	if err != nil {
		log.Println("Gagal membaca sheet:", err)
		return nil, fmt.Errorf("gagal membaca sheet: %w", err)
	}
	defer rows.Close()

	targetHeaders := map[string]bool{
		"Slip Type":       true,
		"Vendor Inv No.":  true,
		"Customer Name":   true,
		"Stmt.No.":        true,
		"Loc Cur":         true,
		"Proxy Y/N":       true,
		"Final Amt":       true,
		"Final Tax Amt":   true,
		"Final Total Amt": true,
		"Acc Ref No.":     true,
	}

	var headers []string
	if rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			log.Println("Gagal mendapatkan kolom dari header:", err)
			return nil, fmt.Errorf("gagal mendapatkan kolom dari header: %w", err)
		}
		headers = cols
	}

	headerIndex := make(map[string]int)
	for i, header := range headers {
		if targetHeaders[header] {
			headerIndex[header] = i
		}
	}

	for target := range targetHeaders {
		if _, exists := headerIndex[target]; !exists {
			log.Printf("Header '%s' tidak ditemukan di file Excel", target)
			return nil, fmt.Errorf("header '%s' tidak ditemukan di file Excel", target)
		}
	}

	var data []map[string]interface{}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			log.Println("Gagal mendapatkan kolom dari row:", err)
			return nil, fmt.Errorf("gagal mendapatkan kolom dari row: %w", err)
		}

		rowData := make(map[string]interface{})
		for header, index := range headerIndex {
			if index < len(row) {
				switch header {
				case "Final Total Amt", "Final Tax Amt", "Final Amt":
					value, err := strconv.ParseFloat(row[index], 64)
					if err != nil {
						log.Printf("Gagal mengonversi nilai '%s' menjadi float64: %v", row[index], err)
						value = 0
					}
					rowData[header] = value
				default:
					rowData[header] = row[index]
				}
			}
		}
		if len(rowData) > 0 {
			data = append(data, rowData)
		}
	}

	return data, nil
}
