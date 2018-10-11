package main

import (
	"fmt"
	"sunrise-sunset-client/client"
	"time"
)

func main() {
	format := "15:04:05"
	date := time.Now()
	c := client.NewSunriseSunsetClient(44.8040, 20.4651, date, "Europe/Belgrade")
	res, _ := c.GetSunriseAndSunsetTimesWithTimezone()
	fmt.Printf("Date: %s\nSunrise: %s\nSunset: %s\nDay length: %s\n", date.Format("2006-01-02"),
		res.Sunrise.Format(format), res.Sunset.Format(format), res.DayLength)
}
