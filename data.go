package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/openweathermap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Current_weather_storage struct {
	Zip      string                            `json:"Zip"`
	Dt       int                               `json:"Dt"`
	OpenData openweathermap.OneCallCurrentData `json:"OpenData"`
}

type Hourly_weather_storage struct {
	Zip      string                           `json:"Zip"`
	Dt       int                              `json:"Dt"`
	OpenData openweathermap.OneCallHourlyData `json:"OpenData"`
}

type Daily_weather_storage struct {
	Zip      string                          `json:"Zip"`
	Dt       int                             `json:"Dt"`
	OpenData openweathermap.OneCallDailyData `json:"OpenData"`
}

type Alert_weather_storage struct {
	Zip      string                          `json:"Zip"`
	Dt       int                             `json:"Dt"`
	OpenData openweathermap.OneCallAlertData `json:"OpenData"`
}

type Time_info struct {
	sunrise string
	sunset  string
	period  string
}

type Event struct {
	Summary  string `json:"Summary"`
	Start    string `json:"Start"`
	End      string `json:"End"`
	Duration string `json:"Duration"`
}

type smart_data struct {
	Current_weather Current_weather_storage  `json:"Current_weather"`
	Hourly_weather  []Hourly_weather_storage `json:"Hourly_weather"`
	Daily_weather   []Daily_weather_storage  `json:"Daily_weather"`
	Alert_weather   []Alert_weather_storage  `json:"Alert_weather"`
	Time_info       Time_info                `json:"Time_info"`
	Events          []Event                  `json:"Events"`
}

// func load_data(mydata *smart_data, loc weather_location) {
func load_data(mydata *smart_data, conf config) {
	getWeather(mydata, conf)
	var local_temp float64
	local_temp, _ = get_temperature(conf)
	mydata.Current_weather.OpenData.Temp = local_temp
	local_humidity, _ := get_humidity(conf)
	mydata.Current_weather.OpenData.Humidity = local_humidity
	// store time info
	mydata.Time_info = get_time_info(conf)
	getCalendar(mydata, conf)
}

func get_humidity(conf config) (int, error) {
	var bearer = "Bearer " + conf.Homeassistant.Token
	url := conf.Homeassistant.Url + "/api/states/" + conf.Homeassistant.Humidity_sensor
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", bearer)
	// add content type header application json
	req.Header.Add("Content-Type", "application/json")
	// do the query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// close the response
	defer resp.Body.Close()
	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// parse the response body
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		panic(err)
	}
	// get the temperature
	m := f.(map[string]interface{})
	// convert string state to int
	result, err := strconv.Atoi(m["state"].(string))
	return result, err
}

// function to retrieve temperature from sensor  in Home Assistant
func get_temperature(conf config) (float64, error) {
	var bearer = "Bearer " + conf.Homeassistant.Token
	url := conf.Homeassistant.Url + "/api/states/" + conf.Homeassistant.Temp_sensor
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", bearer)
	// add content type header application json
	req.Header.Add("Content-Type", "application/json")
	// do the query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// close the response
	defer resp.Body.Close()
	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// parse the response body
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		panic(err)
	}
	// get the temperature
	m := f.(map[string]interface{})
	var result float64
	result, err = strconv.ParseFloat(m["state"].(string), 64)
	// fmt.Println(m["attributes"])
	return result, err
}

func getWeather(mydata *smart_data, conf config) (string, error) {
	w, err := openweathermap.NewOneCall("C", "en", conf.Openweather_api_key, []string{})
	if err != nil {
		log.Fatal(err)
	} else {
		err = w.OneCallByCoordinates(&openweathermap.Coordinates{Latitude: conf.Latitude, Longitude: conf.Longitude})
		if err != nil {
			log.Fatal(err)
		} else {
			mydata.Current_weather = Current_weather_storage{
				Zip:      conf.Zip,
				Dt:       w.Current.Dt,
				OpenData: openweathermap.OneCallCurrentData(w.Current),
			}
			// fill hourly
			for i := 1; i < 6; i++ {
				hour := Hourly_weather_storage{
					Zip:      conf.Zip,
					Dt:       w.Hourly[i].Dt,
					OpenData: openweathermap.OneCallHourlyData(w.Hourly[i]),
				}
				mydata.Hourly_weather = append(mydata.Hourly_weather, hour)
			}
			// fill daily
			for i := 1; i < 6; i++ {
				day := Daily_weather_storage{
					Zip:      conf.Zip,
					Dt:       w.Daily[i].Dt,
					OpenData: openweathermap.OneCallDailyData(w.Daily[i]),
				}
				mydata.Daily_weather = append(mydata.Daily_weather, day)
			}
			// fill alerts
			// set curdt as now in epoch
			curdt := int(time.Now().Unix())
			for i := 0; i < len(w.Alerts); i++ {
				alert := Alert_weather_storage{
					Zip:      conf.Zip,
					Dt:       curdt,
					OpenData: openweathermap.OneCallAlertData(w.Alerts[i]),
				}
				mydata.Alert_weather = append(mydata.Alert_weather, alert)
				fmt.Println(alert)
			}
		}
	}
	return "", nil
}

// reply body : {"results":{"sunrise":"2023-03-19T06:03:25+00:00","sunset":"2023-03-19T18:12:14+00:00","solar_noon":"2023-03-19T12:07:49+00:00","day_length":43729,"civil_twilight_begin":"2023-03-19T05:43:49+00:00","civil_twilight_end":"2023-03-19T18:31:50+00:00","nautical_twilight_begin":"2023-03-19T05:19:49+00:00","nautical_twilight_end":"2023-03-19T18:55:50+00:00","astronomical_twilight_begin":"2023-03-19T04:55:49+00:00","astronomical_twilight_end":"2023-03-19T19:19:50+00:00"},"status":"OK"}
type Response struct {
	Results struct {
		Sunrise                 string `json:"sunrise"`
		Sunset                  string `json:"sunset"`
		SolarNoon               string `json:"solar_noon"`
		DayLength               int    `json:"day_length"`
		CivilTwilightBegin      string `json:"civil_twilight_begin"`
		CivilTwilightEnd        string `json:"civil_twilight_end"`
		NauticalTwilightBegin   string `json:"nautical_twilight_begin"`
		NauticalTwilightEnd     string `json:"nautical_twilight_end"`
		AstronomicalTwilightBeg string `json:"astronomical_twilight_begin"`
		AstronomicalTwilightEnd string `json:"astronomical_twilight_end"`
	} `json:"results"`
	Status string `json:"status"`
}

func get_time_info(conf config) Time_info {
	var res Time_info

	api_url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%f&lng=%f&date=today&formatted=0", conf.Latitude, conf.Longitude)
	// get the json from the api
	req, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("Status code != 200")
	}
	if resp.Body == nil {
		panic("Body is nil")
	} else {
		defer resp.Body.Close()
	}
	// decode the json
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	t1, err := time.Parse(time.RFC3339, response.Results.Sunrise)
	if err != nil {
		panic(err)
	}
	// print t1 as date/time string in New York time zone
	tloc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	res.sunrise = t1.In(tloc).Format("03:04 PM")
	t2, err := time.Parse(time.RFC3339, response.Results.Sunset)
	if err != nil {
		panic(err)
	}
	res.sunset = t2.In(tloc).Format("03:04 PM")
	// period is 'day' if now is between sunrise and sunset
	// period is 'night' if now is not between sunrise and sunset
	now := time.Now().In(tloc)
	if now.After(t1) && now.Before(t2) {
		res.period = "day"
	} else {
		res.period = "night"
	}
	//res.now = now.Format("15:04")
	return res
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getCalendar(mydata *smart_data, conf config) {
	// get the google calendar
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	// get event list from calendar
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, err := getClient(config)
	if err != nil {
		log.Printf("Unable to get client: %v", err)
	} else {
		calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		tStart := time.Now().Format(time.RFC3339)
		tEnd := time.Now().AddDate(0, 0, 14).Format(time.RFC3339)
		events, err := calendarService.Events.List(conf.Gcal_id + "@group.calendar.google.com").
			ShowDeleted(false).SingleEvents(true).TimeMin(tStart).TimeMax(tEnd).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
		}
		nbItems := 0
		if len(events.Items) > 0 {
			for _, item := range events.Items {
				// fmt.Printf("%s %s %s\n", item.Summary, item.Start.DateTime, item.End.DateTime)
				if nbItems == 5 {
					break
				}
				var duration time.Duration
				var startD, endD string
				// if item is all day
				if item.Start.DateTime == "" {
					// get nb of days
					t1, err := time.Parse("2006-01-02", item.Start.Date)
					if err != nil {
						panic(err)
					}
					t2, err := time.Parse("2006-01-02", item.End.Date)
					if err != nil {
						panic(err)
					}
					// startD = t1.Local().Format("Mon 02/01")
					// endD = t2.Local().Format("Mon 02/01")
					startD = t1.Local().Format("2006-01-02")
					endD = t2.Local().Format("2006-01-02")
					duration = t2.Sub(t1)
				} else {
					// print duration
					t1, err := time.Parse(time.RFC3339, item.Start.DateTime)
					if err != nil {
						panic(err)
					}
					t2, err := time.Parse(time.RFC3339, item.End.DateTime)
					if err != nil {
						panic(err)
					}
					duration = t2.Sub(t1)
					// startD = t1.Local().Format("Mon 02/01 15:04")
					// endD = t2.Local().Format("Mon 02/01 15:04")
					startD = t1.Local().Format("2006-01-02 15:04")
					endD = t2.Local().Format("2006-01-02 15:04")
				}
				res := ""
				if duration.Hours() >= 24 {
					res += fmt.Sprintf("%dd", int(duration.Hours())/24)
				}
				if int(duration.Hours())%24 != 0 {
					res += fmt.Sprintf("%dh", int(duration.Hours())%24)
				}
				if int(duration.Minutes())%60 != 0 {
					res += fmt.Sprintf("%dm", int(duration.Minutes())%60)
				}
				// fmt.Printf("Duration: %s\n", res)
				nbItems++
				event := Event{
					Summary:  item.Summary,
					Start:    startD,
					End:      endD,
					Duration: res,
				}
				mydata.Events = append(mydata.Events, event)
				// print event
				// fmt.Printf("%s (%s - %s) %s\n", item.Summary, startD, endD, res)
			}
			// } else {
			// 	// fmt.Println("No upcoming events found.")
		}
	}
}
