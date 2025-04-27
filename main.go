package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
	data "github.com/mnadev/adhango/pkg/data"
	util "github.com/mnadev/adhango/pkg/util"
)

const (
	fajrDegree = 20.0
	ishaDegree = 18.0
	lat        = -7.235465
	lon        = 112.760311
)

var prayersToTime = map[string]time.Time{}

func main() {
	updateAndSchedule()

	for {
		waitUntilMidnight()
		updateAndSchedule()
	}
}

func updateAndSchedule() {
	prayerTimes, err := calculatePrayersTime()
	if err != nil {
		log.Fatalf("error when calculating prayers time: %v", err)
	}

	prayersToTime = map[string]time.Time{
		"Fajr":    prayerTimes.Fajr,
		"Dhuhr":   prayerTimes.Dhuhr,
		"Asr":     prayerTimes.Asr,
		"Maghrib": prayerTimes.Maghrib,
		"Isha":    prayerTimes.Isha,
	}

	for name, t := range prayersToTime {
		fmt.Println("Scheduled", name, t)
		go scheduleNotification(name, t)
	}
}

func waitUntilMidnight() {
	now := time.Now()

	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	durationUntilMidnight := nextMidnight.Sub(now)
	log.Printf("Waiting until midnight, duration: %v", durationUntilMidnight)

	time.Sleep(durationUntilMidnight)
}

func calculatePrayersTime() (*calc.PrayerTimes, error) {
	date := data.NewDateComponents(time.Now())
	params := calc.NewCalculationParametersBuilder().
		SetMadhab(calc.SHAFI_HANBALI_MALIKI).
		SetMethod(calc.OTHER).
		SetFajrAngle(fajrDegree).
		SetIshaAngle(ishaDegree).
		SetMethodAdjustments(calc.PrayerAdjustments{
			DhuhrAdj: 1,
		}).
		Build()

	coords, err := util.NewCoordinates(lat, lon)
	if err != nil {
		fmt.Printf("got error %+v", err)
		return &calc.PrayerTimes{}, err
	}

	prayerTimes, err := calc.NewPrayerTimes(coords, date, params)
	if err != nil {
		fmt.Printf("got error %+v", err)
		return &calc.PrayerTimes{}, err
	}

	err = prayerTimes.SetTimeZone("Asia/Jakarta")
	if err != nil {
		fmt.Printf("got error %+v", err)
		return &calc.PrayerTimes{}, err
	}

	return prayerTimes, nil
}

func scheduleNotification(prayerName string, prayerTime time.Time) {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)
	targetTime := prayerTime.In(loc)

	if targetTime.Before(now) {
		log.Printf("%s prayer already passed", prayerName)
		return
	}

	duration := targetTime.Sub(now)
	log.Printf("Scheduled %s prayer in %v", prayerName, duration)

	time.Sleep(duration)

	err := sendNotification("AdhanPulse", fmt.Sprintf("It's time for %s prayer", prayerName))
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

func sendNotification(title, message string) error {
	cmd := exec.Command("notify-send", "-i", "dialog-information", "-u", "normal", title, message)
	return cmd.Run()
}
