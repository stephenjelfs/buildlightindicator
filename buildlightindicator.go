package buildlightindicator

import (
	"github.com/stephenjelfs/hid"
	"errors"
	"log"
	"fmt"
	"time"
)

const (
	LED_GREEN byte = 1
	LED_RED byte = 2
	LED_BLUE byte = 4
)

func SwitchLightToRed() {
	switchLightTo(LED_RED)
}

func SwitchLightToGreen() {
	switchLightTo(LED_GREEN)
}

func SwitchLightToBlue() {
	switchLightTo(LED_BLUE)
}

func SwitchLightToOff() {
	switchLightTo(0)
}

func switchLightTo(ledColor byte) error {
	device, err := connectToLightDevice(4037, 45184)

	if err != nil {
		log.Println(err)

		log.Println("Known USB devices:")
		for _, devInfo := range(hid.Enumerate(0, 0)) {
			log.Println("VendorID:", devInfo.VendorID, "ProductID:", devInfo.ProductID, "Product", devInfo.Product)
		}

		return err
	}

	defer device.Close()

	setLedsPowerFullOff(device);

	if ledColor != 0 {
		breathEffect(device, LED_BLUE);
		setLedColorFullOn(device, ledColor);
	}

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

func breathEffect(device *hid.Device, ledColor byte) {
	turnLedOn(device, ledColor);

	// fade in
	for i := 0; i <= 30; i++ {
		setLedsPower(device, ledColor, byte(i * 3))
		time.Sleep(50 * time.Millisecond)
	}

	// fade out
	for i := 30; i >= 0; i-- {
		setLedsPower(device, ledColor, byte(i * 3))
		time.Sleep(50 * time.Millisecond)
	}

	turnLedOff(device, ledColor)
}

func setLedColorFullOn(device *hid.Device, ledColor byte) {
	turnLedOn(device, ledColor)
	setLedsPower(device, ledColor, 100)
}

func setLedsPowerFullOff(device *hid.Device) {
	setLedsPower(device, LED_GREEN,0)
	setLedsPower(device, LED_RED,0)
	setLedsPower(device, LED_BLUE,0)

	turnLedOff(device, LED_GREEN)
	turnLedOff(device, LED_RED)
	turnLedOff(device, LED_BLUE)
}

func setLedsPower(device *hid.Device, ledColor byte, power byte) {
	var ordinal byte;

	switch(ledColor) {
		case LED_GREEN:
			ordinal = 0
		case LED_RED:
			ordinal = 1
		case LED_BLUE:
			ordinal = 2
	}

	device.Write([]byte {101, 34, ordinal, power, 0, 0, 0, 0})
}

func turnLedOn(device *hid.Device, ledColor byte) {
	device.Write([]byte {101, 20, ledColor, 0, 0, 0, 0, 0}) // turn off flash
	device.Write([]byte {101, 12, ledColor, 0, 0, 0, 0, 0}) // turn on
}

func turnLedOff(device *hid.Device, ledColor byte) {
	device.Write([]byte {101, 20, 0, ledColor, 0, 0, 0, 0, 0}) // turn off flash
	device.Write([]byte {101, 12, 0, 0, 0, 0, 0, 0}) // turn off
}
