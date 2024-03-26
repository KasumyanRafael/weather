package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const city = "Beslan"
const apiURL = "https://api.openweathermap.org/data/2.5/weather?q=Beslan&appid=YOUR_API_KEY"

func main() {
	file, err := os.OpenFile("weather_data.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for {
		err := saveWeatherData(writer)
		if err != nil {
			fmt.Println("Error saving weather data:", err)
		}

		time.Sleep(1 * time.Hour)
	}
}

func saveWeatherData(writer *csv.Writer) error {
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var data map[string]interface{}
	_, err = io.ReadFull(resp.Body, &data)
	if err != nil {
		return err
	}

	temp := data["main"].(map[string]interface{})["temp"].(float64)
	description := data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)
	timestamp := time.Now().Unix()

	err = writer.Write([]string{city, fmt.Sprintf("%.2f", temp), description, fmt.Sprintf("%d", timestamp)})
	if err != nil {
		return err
	}

	return nil
}
