package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/image/bmp"
)

const IMAGE_WIDTH = 800
const IMAGE_HEIGHT = 480

type weather_location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Zip string  `json:"zip"`
}

type homeassistant_config struct {
	Token           string `json:"homeassistant_token"`
	Url             string `json:"homeassistant_url"`
	Temp_sensor     string `json:"temp_sensor"`
	Humidity_sensor string `json:"humidity_sensor"`
}

type config struct {
	Latitude            float64              `json:"latitude"`
	Longitude           float64              `json:"longitude"`
	Zip                 string               `json:"zip"`
	Homeassistant       homeassistant_config `json:"homeassistant"`
	Openweather_api_key string               `json:"openweather_api_key"`
	Gcal_id             string               `json:"gcal_id"`
}

func main() {
	// load config
	config := new(config)
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), config)

	//home := weather_location{Lat: 41.593147, Lon: -118.550095, Zip: "33175"}
	fonts := new(preloaded_fonts)
	// parse ttf font file from fonts directory
	load_font(fonts)
	p := color.Palette([]color.Color{color.White, color.Black})
	img := image.NewPaletted(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT), p)
	dt := time.Now()
	// draw box for web view around the image
	draw.Draw(img, image.Rect(0, 0, 800, 1), &image.Uniform{color.Black}, image.ZP, draw.Src)
	draw.Draw(img, image.Rect(0, 479, 800, 480), &image.Uniform{color.Black}, image.ZP, draw.Src)
	draw.Draw(img, image.Rect(0, 0, 1, 480), &image.Uniform{color.Black}, image.ZP, draw.Src)
	draw.Draw(img, image.Rect(799, 0, 800, 480), &image.Uniform{color.Black}, image.ZP, draw.Src)
	// add data
	addLabel(img, 520, 30, dt.Format("Monday January 2"), fonts, 25)

	// add week number
	//addLabel(img, 520, 55, "Week "+dt.Format("02"), fonts, 20)

	mydata := new(smart_data)
	load_data(mydata, *config)

	// draw current weather
	draw_current_weather(img, 0, 0, mydata, fonts)
	draw_moon_phase(img, 500, 50, mydata, fonts)
	//draw_hourly_forecast(img, 0, 250, mydata, fonts)
	draw_daily_forecast(img, 00, 300, mydata, fonts)
	// draw_calendar(img, 330, 120, mydata, fonts)
	// list calendar entries
	//draw_calendar(img, 520, 80, mydata, fonts)
	draw_events_icons(img, 700, 280, mydata, fonts)
	// save bmp file
	outfile, err := os.Create("output/smarteink.bmp")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	_ = bmp.Encode(outfile, img)
	// // print smart_data as json
	// b, err := json.Marshal(mydata)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var prettyJson bytes.Buffer
	// json.Indent(&prettyJson, b, "", "  ")
	// println(prettyJson.String())

}
