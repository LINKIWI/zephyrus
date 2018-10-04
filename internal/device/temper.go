package device

import (
	"fmt"
	"time"

	"zephyrus/schemas"

	"github.com/zserge/hid"
)

// TemperDeviceID is the expected hardware device ID for the USB Temper, formatted as
// vendor:product:revision:interface.
const TemperDeviceID = "413d:2107:0000:01"

// TemperDeviceIOTime is a time.Duration describing the maximum allowed I/O time when reading or
// writing from/to the device.
const TemperDeviceIOTime = 1 * time.Second

// TemperClient is a small client library implementing the Sensor interface for interacting with
// a USB-attached Temper device.
type TemperClient struct {
	// HID device backend.
	device hid.Device
	// Unique identifier, set by the user.
	identifier string
	// Current status of the device.
	status schemas.Status
}

// NewTemperClient attempts to find an attached Temper device and creates a client instance with
// the specified identifier name.
func NewTemperClient(identifier string) (*TemperClient, error) {
	var temperDev hid.Device

	hid.UsbWalk(func(dev hid.Device) {
		info := dev.Info()
		walkDeviceID := fmt.Sprintf(
			"%04x:%04x:%04x:%02x",
			info.Vendor,
			info.Product,
			info.Revision,
			info.Interface,
		)

		if walkDeviceID == TemperDeviceID {
			temperDev = dev
		}
	})

	if temperDev == nil {
		return nil, fmt.Errorf("temper: Unable to find USB Temper device")
	}

	return &TemperClient{
		identifier: identifier,
		device:     temperDev,
		status:     schemas.Status_UNKNOWN,
	}, nil
}

// Open opens the located HID device.
func (t *TemperClient) Open() error {
	err := t.device.Open()
	defer func() {
		if err != nil {
			t.status = schemas.Status_ERROR
		} else {
			t.status = schemas.Status_OPENED
		}
	}()

	return err
}

// Close opens the located HID device.
func (t *TemperClient) Close() error {
	t.device.Close()
	t.status = schemas.Status_CLOSED

	return nil
}

// GetIdentifier simply returns the user-set identifier.
// Note that this method does not attempt to read any data from the actual device.
func (t *TemperClient) GetIdentifier() (string, error) {
	return t.identifier, nil
}

// GetStatus reports the current device state.
func (t *TemperClient) GetStatus() schemas.Status {
	return t.status
}

// GetTemperature requests a temperature reading from the device and returns it as a float64 in
// celsius units.
func (t *TemperClient) GetTemperature() (float64, error) {
	readTempCmd := []byte{0x01, 0x80, 0x33, 0x01, 0x00, 0x00, 0x00, 0x00}

	if _, err := t.device.Write(readTempCmd, TemperDeviceIOTime); err != nil {
		return 0.0, fmt.Errorf("temper: %v", err)
	}

	buf, err := t.device.Read(-1, TemperDeviceIOTime)
	if err != nil {
		return 0.0, fmt.Errorf("temper: %v", err)
	}

	temperature := (256*float64(buf[2]) + float64(buf[3])) / 100.0 // Celsius units
	return temperature, nil
}
