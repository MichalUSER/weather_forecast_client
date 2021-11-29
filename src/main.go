package main

import (
	"fmt"
	"time"

	"bytes"
	"strconv"

	"encoding/json"
	"net/http"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
)

var firmataAdaptor *firmata.Adaptor

var temps []float64
var tempsSum float64

func average() string {
	return fmt.Sprintf("%.2f", tempsSum/float64(len(temps)))
}

func clearTemps() {
	temps = nil
	tempsSum = 0
}

func measureTemp() {
	val, err := firmataAdaptor.AnalogRead("0")
	if err != nil {
		fmt.Println(err)
		return
	}

	voltage := (float64(val) * 5) / 1024
	//fmt.Println("voltage:", voltage)
	tempC := (voltage - 0.5) * 100
	//tempC := ((float64(val) * 0.1039) - 50.0)
	fmt.Println("temp:", tempC)

	tempCf, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", tempC), 64)
	temps = append(temps, tempCf)
	tempsSum += tempCf
}

func newTemp() *Temp {
	t := time.Now()
	return &Temp{
		Y:           t.Year(),
		M:           int(t.Month()),
		D:           t.Day(),
		H:           t.Hour(),
		AverageTemp: average(),
	}
}

func addTemp(temp *Temp) {
	tempJSON, err := json.Marshal(temp)
	if err != nil {
		fmt.Println(err)
	}
	_, err = http.Post("http://192.168.0.110:8080/add_temp", "application/json",
		bytes.NewBuffer(tempJSON))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	firmataAdaptor = firmata.NewAdaptor("/dev/ttyACM0")
	work := func() {
		time.Sleep(1 * time.Second)
		for i := 0; i <= 5; i++ {
			measureTemp()
			time.Sleep(1 * time.Second)
		}
		addTemp(newTemp())
		clearTemps()
		gobot.Every(60*time.Minute, func() {
			for i := 0; i <= 5; i++ {
				measureTemp()
				time.Sleep(1 * time.Second)
			}
			addTemp(newTemp())
			clearTemps()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{},
		work,
	)

	robot.Start()
}
