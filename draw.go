package main

import (
	"fmt"
	"image"
	"math"
	"strconv"
	"time"
	"unicode"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

func addLabel(img *image.Paletted, x, y int, label string, fonts *preloaded_fonts, size float64) {
	pt := freetype.Pt(x, y)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fonts.text)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingFull)
	_, _ = c.DrawString(label, pt)
}

func addIcon(img *image.Paletted, x, y int, label string, fonts *preloaded_fonts, size float64) {
	pt := freetype.Pt(x, y)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fonts.icons)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)
	_, _ = c.DrawString(label, pt)
}

func addCustomIcon(img *image.Paletted, x, y int, label string, fonts *preloaded_fonts, size float64) {
	pt := freetype.Pt(x, y)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fonts.custom)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)
	_, _ = c.DrawString(label, pt)
}

func draw_current_weather(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
	id := strconv.Itoa(data.Current_weather.OpenData.Weather[0].ID)
	iconstring := WEATHER_ICONS[id][data.Time_info.period]
	addIcon(img, x+20, y+140, iconstring, fonts, 120)
	addLabel(img, x+200, y+60, fmt.Sprintf("%.1f", data.Current_weather.OpenData.Temp)+" °C", fonts, 40)
	// add temperature in Fahrenheit
	//addLabel(img, x+130, y+80, fmt.Sprintf("%.1f", data.Current_weather.OpenData.Temp*9/5+32)+" °F", fonts, 30)
	addIcon(img, x+200, y+115, "\uf00d", fonts, 25)
	addLabel(img, x+230, y+115, fmt.Sprintf("UVI %.0f", data.Current_weather.OpenData.UVI), fonts, 25)
	r := []rune(data.Current_weather.OpenData.Weather[0].Description)
	r[0] = unicode.ToUpper(r[0])
	addLabel(img, x+20, y+195, string(r), fonts, 18)

	addIcon(img, x+20, y+225, "\uf07a", fonts, 20)
	addLabel(img, x+50, y+225, strconv.Itoa(data.Current_weather.OpenData.Humidity)+" %", fonts, 18)
	addIcon(img, x+150, y+225, "\uf050", fonts, 20)
	addLabel(img, x+180, y+225, fmt.Sprintf("%.1f", data.Current_weather.OpenData.WindSpeed*3.6)+" km/h", fonts, 18)
	// addIcon(img, x+120, y+150, "\uf079", fonts, 20)
	// addLabel(img, x+150, y+150, strconv.Itoa(data.Current_weather.OpenData.Pressure)+" hPa", fonts, 18)
	// add sunrise and sunset
	addIcon(img, x+20, y+255, "\uf051", fonts, 20)
	addLabel(img, x+50, y+255, data.Time_info.sunrise, fonts, 18)
	addIcon(img, x+150, y+255, "\uf052", fonts, 20)
	addLabel(img, x+180, y+255, data.Time_info.sunset, fonts, 18)
}

// large
func draw_moon_phase(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
	moon := NewMoon(time.Now())

	moonString := int('\uf0d0') - 1 + int(moon.Phase()*27)
	if int(moon.Phase()*27) == 0 {
		moonString = int('\uf0EB')
	}
	addIcon(img, x+100, y+90, string(moonString), fonts, 80)
	nextFt := moon.FullMoon() // next full moon if we are in the first half of the moon cycle
	dir := "\uf057"           // up
	if moon.Phase() > 0.5 {
		dir = "\uf088"               // down
		nextFt = moon.NextFullMoon() // next full moon if we are in the second half of the moon cycle
	}
	addIcon(img, x+40, y+140, dir, fonts, 35)
	addLabel(img, x+60, y+135, moon.PhaseName(), fonts, 18)
	// add next full moon
	sec, dec := math.Modf(nextFt)
	nextF := time.Unix(int64(sec), int64(dec*(1e9)))
	addIcon(img, x+40, y+175, "\uf0dd", fonts, 20)
	addLabel(img, x+60, y+175, nextF.Format("Jan 02 03:04 PM"), fonts, 18)
}

// compact
// func draw_moon_phase(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
// 	moon := NewMoon(time.Now())
// 	moonString := int('\uf0d0') - 1 + int(moon.Phase()*27)
// 	addIcon(img, x, y+90, string(moonString), fonts, 80)
// 	nextFt := moon.FullMoon() // next full moon if we are in the first half of the moon cycle
// 	dir := "\uf057"           // up
// 	if moon.Phase() > 0.5 {
// 		dir = "\uf088"               // down
// 		nextFt = moon.NextFullMoon() // next full moon if we are in the second half of the moon cycle
// 	}
// 	addIcon(img, x+70, y+60, dir, fonts, 35)
// 	addLabel(img, x+90, y+55, moon.PhaseName(), fonts, 18)
// 	// add next full moon
// 	sec, dec := math.Modf(nextFt)
// 	nextF := time.Unix(int64(sec), int64(dec*(1e9)))
// 	addIcon(img, x+70, y+80, "\uf0dd", fonts, 20)
// 	addLabel(img, x+90, y+80, nextF.Format("Jan 02 03:04 PM"), fonts, 16)
// }

// draw_hourly_forecast(img, 0, 300, mydata, fonts)
// func draw_hourly_forecast(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
// 	// draw the next 5 hours
// 	for i := 0; i < 5; i++ {
// 		hourS := time.Unix(int64(data.Hourly_weather[i].OpenData.Dt), 0).Format("15H")
// 		addLabel(img, x+20+i*50, y+25, hourS, fonts, 15)
// 		id := strconv.Itoa(data.Hourly_weather[i].OpenData.Weather[0].ID)
// 		iconstring := WEATHER_ICONS[id][data.Time_info.period]
// 		addIcon(img, x+20+i*50, y+50, iconstring, fonts, 20)
// 		addLabel(img, x+20+i*50, y+70, fmt.Sprintf("%.0f", data.Hourly_weather[i].OpenData.Temp)+" °C", fonts, 15)
// 		addLabel(img, x+20+i*50, y+85, fmt.Sprintf("%.0f", data.Hourly_weather[i].OpenData.Pop*100)+" %", fonts, 15)
// 	}
// }

func draw_daily_forecast(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
	// draw the next 5 days
	offset := 160
	for i := 0; i < 5; i++ {
		dayS := time.Unix(int64(data.Daily_weather[i].OpenData.Dt), 0).Format("Monday")
		id := strconv.Itoa(data.Daily_weather[i].OpenData.Weather[0].ID)
		iconstring := WEATHER_ICONS[id][data.Time_info.period]
		addLabel(img, x+20+i*offset, y+25, dayS, fonts, 18)
		addIcon(img, x+40+i*offset, y+90, iconstring, fonts, 50)
		// add description
		r := []rune(data.Daily_weather[i].OpenData.Weather[0].Description)
		r[0] = unicode.ToUpper(r[0])
		addLabel(img, x+20+i*offset, y+120, string(r), fonts, 16)
		addIcon(img, x+20+i*offset, y+150, "\uf053", fonts, 20)
		addLabel(img, x+35+i*offset, y+150, fmt.Sprintf("%.0f", data.Daily_weather[i].OpenData.Temp.Min)+"°C", fonts, 18)
		addIcon(img, x+85+i*offset, y+150, "\uf055", fonts, 20)
		addLabel(img, x+100+i*offset, y+150, fmt.Sprintf("%.0f", data.Daily_weather[i].OpenData.Temp.Max)+"°C", fonts, 18)
		// temperature again, but in Fahrenheit
		// addIcon(img, x+20+i*offset, y+180, "\uf053", fonts, 20)
		// addLabel(img, x+35+i*offset, y+180, fmt.Sprintf("%.0f", data.Daily_weather[i].OpenData.Temp.Min*1.8+32)+"°F", fonts, 18)
		// addIcon(img, x+80+i*offset, y+180, "\uf055", fonts, 20)
		// addLabel(img, x+95+i*offset, y+180, fmt.Sprintf("%.0f", data.Daily_weather[i].OpenData.Temp.Max*1.8+32)+"°F", fonts, 18)
	}
}

//	func draw_calendar(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
//		// print to stdout event list
//		for i := 0; i < 5; i++ {
//			addLabel(img, x+20, y+20+i*20, data.Events[i].Summary, fonts, 15)
//			addLabel(img, x+20, y+35+i*20, data.Events[i].Start, fonts, 15)
//		}
//	}
// func draw_calendar(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
// 	// print to stdout event list
// 	// loop on all events
// 	for i := 0; i < len(data.Events); i++ {
// 		// add start date first 3 chars
// 		addLabel(img, x, y+20+i*20, data.Events[i].Start[:3], fonts, 15)
// 		addLabel(img, x+50, y+20+i*20, data.Events[i].Start[4:], fonts, 15)
// 		addLabel(img, x+160, y+20+i*20, data.Events[i].Summary+" ("+data.Events[i].Duration+")", fonts, 15)
// 	}
// }

func draw_events_icons(img *image.Paletted, x, y int, data *smart_data, fonts *preloaded_fonts) {
	// loop on all events for today
	for i := 0; i < len(data.Events); i++ {
		// check if event is today
		if data.Events[i].End[:10] == time.Now().Format("2006-01-02") {
			// check if event title is a key in CUSTOM_ICONS
			if _, ok := CUSTOM_ICONS[data.Events[i].Summary]; ok {
				// add icon
				addCustomIcon(img, x-45*i, y, CUSTOM_ICONS[data.Events[i].Summary], fonts, 45)
			}
		}
	}
}
