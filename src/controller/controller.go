package controller

import (
	get "Yosyos/src/api/getlist"
	save "Yosyos/src/api/save"
	A "Yosyos/src/api/test"
	"log"

	"bufio"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	cmd "Yosyos/src/config"
)

func RunController() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	RunPort := viper.GetString("port")

	v1 := router.Group("/api/v1")
	v1.Use(A.PassAuthMiddleware())

	v1.GET("/common/excel/list", get.GetListExcelValue)
	v1.GET("/common/excel/save", save.ProcessData)

	serverReady := make(chan bool)

	go func() {
		log.Printf("Server berjalan di port :%s...", RunPort)
		serverReady <- true
		if err := router.Run(":" + RunPort); err != nil {
			log.Fatalf("Gagal menjalankan server: %v", err)
		}
	}()

	<-serverReady

	fmt.Println("=== Aplikasi Proses Lookup Data ===")
	fmt.Println("Pilih opsi:")
	fmt.Println("1. LookUp Data")
	fmt.Println("2. Exit")

	for {
		fmt.Print("Masukkan pilihan (1/2): ")

		var input string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}

		switch input {
		case "1":
			fmt.Println("Melakukan GET request ke localhost:1111/api/v1/common/excel/save...")
			err := cmd.PerformGetRequest()
			if err != nil {
				fmt.Println("Terjadi kesalahan: ", err)
				fmt.Println("Kembali ke menu utama...\n")
			} else {
				fmt.Println("Proses selesai...")
				fmt.Println("Kembali Ke menu Utama...")
			}
		case "2":
			fmt.Println("Kamsahamnida karna kamu sudah menggunakan BOT Excel, Byee...")
			return
		default:
			fmt.Println("Input tidak valid. Program akan keluar.")
			return
		}
	}
}
