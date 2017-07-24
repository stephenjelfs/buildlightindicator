package main

import (
	"github.com/stephenjelfs/buildlightindicator/hidlight"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"time"
)

type CurrentStatus struct {
	LastSuccess   		string `json:"lastSuccess"`
	LastCompleted 		string `json:"lastCompleted"`
	LastCompletedError	string `json:"lastCompletedError"`
	Running       		string `json:"running"`
	Pending       		string `json:"pending"`
}

type request struct{
	updateOrEmpty	string
	done 			chan CurrentStatus
}

func main() {
	requests := make(chan request)

	go syncHidlight(requests)
	startServer(8080, requests)
}

func syncHidlight(requests chan request) {
	updateComplete := make(chan hidlight.Status)

	pending := ""
	running := ""
	lastSuccess := ""
	lastCompleted := ""
	lastCompletedError := ""

	for {
		select {
			case request := <- requests:
				if request.updateOrEmpty != "" {
					if running == "" {
						updateHidlight(request.updateOrEmpty, updateComplete)
						running = request.updateOrEmpty
					} else {
						pending = request.updateOrEmpty
					}
				}
				request.done <- CurrentStatus {lastSuccess, lastCompleted, lastCompletedError, running, pending}

			case status := <- updateComplete:
				if pending != "" {
					updateHidlight(pending, updateComplete)
					running = pending
					pending = ""
				} else {
					running = ""
				}

				if status.Error != nil {
					lastCompletedError = status.Error.Error()
				} else {
					lastSuccess = status.Color
				}
				lastCompleted = status.Color
		}
	}
}

func updateHidlight(color string, complete chan hidlight.Status) {
	go func() {
		log.Println("Switching light to: " + color)
		err := hidlight.SwitchTo(color)
		if err != nil {
			log.Println(err)
		}
		log.Println("Finished switching light to: " + color)
		complete <- hidlight.Status{color, err}
	}()
}

func startServer(port int, requests chan request) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,"Build light indicator running.")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		done := make(chan CurrentStatus)
		updateOrEmpty := ""

		if r.Method == http.MethodPost {
			updateOrEmpty = r.URL.Query().Get("color")

		}

		requests <- request{updateOrEmpty, done}
		json.NewEncoder(w).Encode(<- done)
	})

	log.Println("Starting server on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}