package main

import (
	"github.com/stephenjelfs/buildlightindicator/hidlight"
	"time"
	"math/rand"
)

func main() {
	hidLight := hidlight.New()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	// TODO: Handle system interrupts
	defer func() {
		hidLight.Shutdown()
		time.Sleep(time.Duration(r1.Int31n(1000)) * time.Millisecond)
	}()

	go func() {
		hidLight.SwitchOff()
		time.Sleep(time.Duration(r1.Int31n(1000)) * time.Millisecond)
	}()

	go func() {
		hidLight.SwitchToRed()
		time.Sleep(time.Duration(r1.Int31n(1000)) * time.Millisecond)
	}()

	go func() {
		hidLight.SwitchToGreen()
		time.Sleep(time.Duration(r1.Int31n(1000)) * time.Millisecond)
	}()

	time.Sleep(time.Duration(r1.Int31n(1000)) * time.Millisecond)
}
