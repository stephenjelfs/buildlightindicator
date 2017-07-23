package main

import (
	"github.com/stephenjelfs/buildlightindicator/hidlight"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
)

type JsonStatus struct {
	Color    string `json:"string"`
	ErrorMsg string `json:"errorMsg"`
}

func main() {
	hidLight := hidlight.New()

	reportedStatus := make(chan hidlight.Status)
	desiredColor := make(chan string)

	go syncHidLight(&hidLight, reportedStatus, desiredColor)
	startServer(8080, reportedStatus, desiredColor)
}

func syncHidLight(hidLight *hidlight.Controller, reportedStatus chan hidlight.Status, desiredColor chan string) {
	lastDesiredColor := "blue"
	updateLight := true

	go func() {
		for {
			lastDesiredColor = <- desiredColor
			updateLight = true
		}
	}()

	for {
		if updateLight {
			updateLight = false

			switch lastDesiredColor {
			case "red":
				hidLight.SwitchToRed(reportedStatus)
			case "blue":
				hidLight.SwitchToBlue(reportedStatus)
			case "green":
				hidLight.SwitchToGreen(reportedStatus)
			}
		}
	}
}

func startServer(port int, reportedStatus chan hidlight.Status, desiredColor chan string) {
	lastReportedStatus := hidlight.Status{}

	go func() {
		for {
			lastReportedStatus = <- reportedStatus
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,"Build light indicator running.")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		errorMsg := ""

		if (lastReportedStatus.Error != nil) {
			errorMsg = lastReportedStatus.Error.Error()
		}

		json.NewEncoder(w).Encode(JsonStatus{lastReportedStatus.Color,errorMsg})
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		desiredColor <- r.URL.Query().Get("color")
	})

	log.Println("Starting server on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}