package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

type SunriseSunsetClient struct {
	Latitude  float64
	Longitude float64
	Date      time.Time
	Timezone  string
}

func NewSunriseSunsetClient(lat, lng float64, date time.Time, timezone string) (*SunriseSunsetClient) {
	return &SunriseSunsetClient{
		Latitude:  lat,
		Longitude: lng,
		Date:      date,
		Timezone:  timezone,
	}
}

type Response struct {
	Results ResultUTC `json:"results"`
	Status  string    `json:"status"`
}

type ResultUTC struct {
	Sunrise                   time.Time `json:"sunrise"`
	Sunset                    time.Time `json:"sunset"`
	SolarNoon                 time.Time `json:"solar_noon"`
	DayLength                 int       `json:"day_length"`
	CivilTwilightBegin        time.Time `json:"civil_twilight_begin"`
	CivilTwilightEnd          time.Time `json:"civil_twilight_end"`
	NauticalTwilightBegin     time.Time `json:"nautical_twilight_begin"`
	NauticalTwilightEnd       time.Time `json:"nautical_twilight_end"`
	AstronomicalTwilightBegin time.Time `json:"astronomical_twilight_begin"`
	AstronomicalTwilightEnd   time.Time `json:"astronomical_twilight_end"`
}

type ResultWithTimezone struct {
	Sunrise                   time.Time
	Sunset                    time.Time
	SolarNoon                 time.Time
	DayLength                 string
	CivilTwilightBegin        time.Time
	CivilTwilightEnd          time.Time
	NauticalTwilightBegin     time.Time
	NauticalTwilightEnd       time.Time
	AstronomicalTwilightBegin time.Time
	AstronomicalTwilightEnd   time.Time
}

func (c *SunriseSunsetClient) GetSunriseAndSunsetTimesWithUTC() (*Response, error) {
	lat := c.Latitude
	lng := c.Longitude
	date := c.Date.Format("2006-01-02")
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%.4f&lng=%.4f&date=%s&formatted=0", lat, lng, date)

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	response := &Response{}
	if err := c.jsonDecode(resp.Body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *SunriseSunsetClient) GetSunriseAndSunsetTimesWithTimezone() (*ResultWithTimezone, error) {
	location, _ := time.LoadLocation(c.Timezone)

	response, err := c.GetSunriseAndSunsetTimesWithUTC()
	if err != nil {
		return nil, err
	}
	return &ResultWithTimezone{
		response.Results.Sunrise.In(location),
		response.Results.Sunset.In(location),
		response.Results.SolarNoon.In(location),
		c.secondsToHuman(response.Results.DayLength),
		response.Results.CivilTwilightBegin.In(location),
		response.Results.CivilTwilightEnd.In(location),
		response.Results.NauticalTwilightBegin.In(location),
		response.Results.NauticalTwilightEnd.In(location),
		response.Results.AstronomicalTwilightBegin.In(location),
		response.Results.AstronomicalTwilightEnd.In(location),
	}, nil
}

func (c *SunriseSunsetClient) secondsToHuman(input int) (result string) {
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if days > 0 {
		result = c.plural(int(days), "day") + c.plural(int(hours), "hour") + c.plural(int(minutes), "minute") + c.plural(int(seconds), "second")
	} else if hours > 0 {
		result = c.plural(int(hours), "hour") + c.plural(int(minutes), "minute") + c.plural(int(seconds), "second")
	} else if minutes > 0 {
		result = c.plural(int(minutes), "minute") + c.plural(int(seconds), "second")
	} else {
		result = c.plural(int(seconds), "second")
	}
	return
}

func (c *SunriseSunsetClient) plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func (c *SunriseSunsetClient) jsonDecode(r io.Reader, v interface{}) (error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
