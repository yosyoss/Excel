package config

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func PerformGetRequest() error {
	time.Sleep(1 * time.Second)

	url := "http://localhost:1111/api/v1/common/excel/save"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Gagal membuat request: %v", err)
	}

	fmt.Print("Masukkan Password: ")
	var auth string
	fmt.Scanln(&auth)
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Gagal melakukan GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Gagal membaca respons: %v", err)
	}

	fmt.Printf("Respons dari server: %s\n", string(body))

	return nil
}
