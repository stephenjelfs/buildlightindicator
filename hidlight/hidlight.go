package hidlight

import (
	"github.com/stephenjelfs/hid"
	"errors"
	"log"
	"fmt"
	"time"
)

type Status struct {
	Color string
	Error error
}

type request struct {
	command HidCommand
	done    chan Status
}

type Controller struct {
	queue chan request
}

func (hidLight Controller) SwitchToRed(done chan Status) {
	hidLight.queue <- request{red(), done}
}

func (hidLight Controller) SwitchToGreen(done chan Status) {
	hidLight.queue <- request{green(), done}
}

func (hidLight Controller) SwitchToBlue(done chan Status) {
	hidLight.queue <- request{blue(), done}
}

func (hidLight Controller) SwitchOff(done chan Status) {
	hidLight.queue <- request{off(), done}
}

func New() Controller {
	queue := make(chan request)

	go func() {
		for {
			handleRequest(<- queue)
		}
	}()

	return Controller{queue}
}

func handleRequest(request request) {
	log.Println("Switching light to:", request.command.name())
	err := runCommandOnDevice(request.command)
	if err != nil {
		log.Println(err)
	}
	log.Println("Finished switching light to:", request.command.name())
	request.done <- Status{request.command.name(), err }
}

type HidCommand interface {
	name() string
	apply(device *hid.Device)
}

type HidLed struct {
	color   string
	code    byte
	ordinal byte
}

func red() HidLed {
	return HidLed{"red",1, 0}
}

func green() HidLed {
	return HidLed{"green",4, 2}
}

func blue() HidLed {
	return HidLed{"blue", 2, 1}
}

func (led HidLed) name() string {
	return led.color
}

type HidOff struct {}

func off() HidOff {
	return HidOff{}
}

func (off HidOff) name() string {
	return "off"
}

func (off HidOff) apply(device *hid.Device) {
	red().turnOffNoPower(device)
	green().turnOffNoPower(device)
	blue().turnOffNoPower(device)
}

func (led HidLed) apply(device *hid.Device) {
	// turn all off
	off().apply(device)
	// breath effect
	blue().breathEffect(device)
	// turn on single led
	led.turnOnFullPower(device)
}

func (led HidLed) breathEffect(device *hid.Device) {
	led.turnOn(device)

	// fade in
	for i := 0; i <= 30; i++ {
		led.setPower(device, byte(i * 3))
		time.Sleep(50 * time.Millisecond)
	}

	// fade out
	for i := 30; i >= 0; i-- {
		led.setPower(device, byte(i * 3))
		time.Sleep(50 * time.Millisecond)
	}

	led.turnOffNoPower(device)
}

func runCommandOnDevice(command HidCommand) error {
	if !hid.Supported() {
		return hid.ErrUnsupportedPlatform
	}

	device, err := connectToLightDevice(4037, 45184)

	if err != nil {
		return err
	}

	defer device.Close()
	command.apply(device)
	return nil
}

func connectToLightDevice(vendorID uint16, productID uint16) (*hid.Device, error) {
	devInfo := hid.Enumerate(vendorID, productID)

	if len(devInfo) == 0 {
		return nil, errors.New(fmt.Sprint("Light device not found: ", "VendorID:", vendorID, "ProductID:", productID))
	}

	device, err := devInfo[0].Open()

	if err != nil {
		log.Println(err, devInfo[0].Path)
		return nil, errors.New(fmt.Sprint("Light device found, but failed to open: ", devInfo[0].Product, ", ", devInfo[0].Manufacturer))
	}

	return device, nil
}

func (led HidLed) turnOnFullPower(device *hid.Device) {
	led.turnOn(device)
	led.setPower(device, 100)
}

func (led HidLed) turnOffNoPower(device *hid.Device) {
	led.turnOff(device)
	led.setPower(device, 0)
}

func (led HidLed) setPower(device *hid.Device, power byte) {
	device.Write([]byte {101, 34, led.ordinal, power, 0, 0, 0, 0})
}

func (led HidLed) turnOn(device *hid.Device) {
	device.Write([]byte {101, 20, led.code, 0, 0, 0, 0, 0}) // turn off flash
	device.Write([]byte {101, 12, led.code, 0, 0, 0, 0, 0}) // turn on
}

func (led HidLed) turnOff(device *hid.Device) {
	device.Write([]byte {101, 20, led.code, 0, 0, 0, 0, 0})    // turn off flash
	device.Write([]byte {101, 12, 0, led.code, 0, 0, 0, 0, 0}) // turn off
}
