package excel

import (
	"log"

	"Yosyos/src/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetListExcelValue(c *gin.Context) {
	var txt model.Parameter
	if errq := c.BindJSON(&txt); errq != nil {
		log.Fatal("Gagal menyimpan filel:", errq)
		c.JSON(400, gin.H{
			"error": "Gagal menyimpan file",
		})
	}
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

func ExcelGet(c *gin.Context) (interface{}, error) {
	f, err := excelize.OpenFile("src/DataExcel/Test2.xlsx")
	if err != nil {
		log.Fatal("Gagal membuka file Excel:", err)
		c.JSON(500, gin.H{
			"error": "Gagal membuka file Excel",
		})
		return nil, err
	}

	rows, err := f.Rows("Sheet1")
	if err != nil {
		log.Fatal("Gagal membaca sheet:", err)
		c.JSON(500, gin.H{
			"error": "Gagal membaca sheet",
		})
		return nil, err
	}

	var headers []string
	if rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			log.Fatal("Gagal mendapatkan kolom dari header:", err)
			c.JSON(500, gin.H{
				"error": "Gagal mendapatkan kolom dari header",
			})
			return nil, err
		}
		headers = cols
	}

	var data []map[string]interface{}

	for rows.Next() {
		row := make(map[string]interface{})

		cols, err := rows.Columns()
		if err != nil {
			log.Fatal("Gagal mendapatkan kolom dari row:", err)
			c.JSON(500, gin.H{
				"error": "Gagal mendapatkan kolom dari row",
			})
			return nil, err
		}

		for i, col := range cols {
			if i < len(headers) {
				row[headers[i]] = col
			}
		}

		data = append(data, row)
	}

	c.JSON(200, gin.H{
		"data": data,
	})
	return data, nil
}
