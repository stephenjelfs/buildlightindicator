package buildlightindicator

import (
	"github.com/stephenjelfs/hid"
	"errors"
	"log"
	"fmt"
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

func switchLightTo(ledColor byte) error {
	//_, err := connectToLightDevice(4037, 45184)
	device, err := connectToLightDevice(1386, 830)

	if err != nil {
		log.Println(err)

		log.Println("Known USB devices:")
		for _, devInfo := range(hid.Enumerate(0, 0)) {
			log.Println("VendorID:", devInfo.VendorID, "ProductID:", devInfo.ProductID, "Product", devInfo.Product)
		}

		return err
	}
	defer device.Close()

	fmt.Println(device)

	return nil
}

func connectToLightDevice(vendorID uint16, productID uint16) (*hid.Device, error) {
	devInfo := hid.Enumerate(vendorID, productID)

	if len(devInfo) == 0 {
		return nil, errors.New(fmt.Sprint("Light device not found: ", "VendorID:", vendorID, "ProductID:", productID))
	}

	device, err := devInfo[0].Open()

	if err != nil {
		return nil, errors.New(fmt.Sprint("Light device found, but failed to open: ", devInfo[0].Product, ", ", devInfo[0].Manufacturer))
	}

	return device, nil
}

