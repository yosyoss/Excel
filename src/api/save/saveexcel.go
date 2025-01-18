package save

import (
	"Yosyos/src/api/getlist"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
		finalAmt, _ := item["Final Amt"].(float64)
		finalTaxAmt, _ := item["Final Tax Amt"].(float64)
		finalTotalAmt, _ := item["Final Total Amt"].(float64)

		data = append(data, model.TaxData{
			SlipType:      item["Slip Type"].(string),
			VendorInvNo:   item["Vendor Inv No."].(string),
			CustomerName:  item["Customer Name"].(string),
			StmtNo:        item["Stmt.No."].(string),
			LocCur:        item["Loc Cur"].(string),
			Proxy:         item["Proxy Y/N"].(string),
			FinalAmt:      finalAmt,
			FinalTaxAmt:   finalTaxAmt,
			FinalTotalAmt: finalTotalAmt,
			AccRefNo:      item["Acc Ref No."].(string),
		})
	}

	groupedData := make(map[string]model.TaxData)

	for _, entry := range data {
		key := fmt.Sprintf("%s-%s", entry.Proxy, entry.StmtNo)

		groupedData[key] = model.TaxData{
			SlipType:      entry.SlipType,
			VendorInvNo:   entry.VendorInvNo,
			CustomerName:  entry.CustomerName,
			StmtNo:        entry.StmtNo,
			LocCur:        entry.LocCur,
			AccRefNo:      entry.AccRefNo,
			Proxy:         entry.Proxy,
			FinalAmt:      groupedData[key].FinalAmt + entry.FinalAmt,
			FinalTaxAmt:   groupedData[key].FinalTaxAmt + entry.FinalTaxAmt,
			FinalTotalAmt: groupedData[key].FinalTotalAmt + entry.FinalTotalAmt,
		}
	}

	var result []model.Response
	for _, entry := range groupedData {
		finalAmtStr := strconv.FormatFloat(entry.FinalAmt, 'f', 2, 64)
		finalTaxAmtStr := strconv.FormatFloat(entry.FinalTaxAmt, 'f', 2, 64)
		finalTotalAmtStr := strconv.FormatFloat(entry.FinalTotalAmt, 'f', 2, 64)

		result = append(result, model.Response{
			SlipType:      entry.SlipType,
			VendorInvNo:   entry.VendorInvNo,
			CustomerName:  entry.CustomerName,
			StmtNo:        entry.StmtNo,
			LocCur:        entry.LocCur,
			AccRefNo:      entry.AccRefNo,
			FinalAmt:      finalAmtStr,
			FinalTaxAmt:   finalTaxAmtStr,
			FinalTotalAmt: finalTotalAmtStr,
		})
	}

	// Simpan ke Excel
	err = saveToExcel(result)
	if err != nil {
		log.Println("Gagal menyimpan ke Excel:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file Excel"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func saveToExcel(data []model.Response) error {
	folderPath := "src/DataExcel/Result"
	fileName := "OutputLookup.xlsx"
	filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
	paymentDate := time.Now().AddDate(0, 0, 1)

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("gagal membuat folder: %w", err)
	}

	f := excelize.NewFile()

	f.NewSheet("Sheet1")

	headers := []string{"No", "Office", "Account Code", "Business Type", "Payment", "Payment date", "Operating Month", "Vendor Inv No.", "Vendor name (User)", "Vendor Inv Date", "STP", "Cur", "Gross", "Tax", "Amount ", "Acc Ref No", "Remark", "Date Payment "}

	for col, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+col)))
		f.SetCellValue("Sheet1", cell, header)
	}

	for i, item := range data {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), "Indonesia")
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), item.SlipType)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), item.BusinessType)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), item.Payment)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), paymentDate.Format("2006-01-02"))
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), item.OperatingMonth)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), item.VendorInvNo)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), item.CustomerName)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", row), item.VendorInvDate)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", row), item.StmtNo)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", row), item.LocCur)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", row), item.FinalAmt)
		f.SetCellValue("Sheet1", fmt.Sprintf("N%d", row), item.FinalTaxAmt)
		f.SetCellValue("Sheet1", fmt.Sprintf("O%d", row), item.FinalTotalAmt)
		f.SetCellValue("Sheet1", fmt.Sprintf("P%d", row), item.AccRefNo)
		f.SetCellValue("Sheet1", fmt.Sprintf("Q%d", row), item.Remark)
		f.SetCellValue("Sheet1", fmt.Sprintf("R%d", row), item.DatePayment)
	}

	err = f.SaveAs(filePath)
	if err != nil {
		return fmt.Errorf("gagal menyimpan file Excel: %w", err)
	}

	log.Printf("File berhasil disimpan ke: %s", filePath)
	return nil
}
