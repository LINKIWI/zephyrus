package device

import (
	"fmt"
	"sync"
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
	// Mutex used to synchronize access to the device.
	mutex sync.Mutex
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
	// Black magic
	readTempCmd := []byte{0x01, 0x80, 0x33, 0x01, 0x00, 0x00, 0x00, 0x00}

	resp, err := t.txRx(readTempCmd, TemperDeviceIOTime)
	if err != nil {
		return 0.0, fmt.Errorf("temper: %v", err)
	}

	temperature := (256*float64(resp[2]) + float64(resp[3])) / 100.0 // Celsius units
	return temperature, nil
}

// Write a byte sequence followed by reading from the device until EOF.
func (t *TemperClient) txRx(writeCmd []byte, timeout time.Duration) ([]byte, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, err := t.device.Write(writeCmd, timeout); err != nil {
		return nil, err
	}

	buf, err := t.device.Read(-1, timeout)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
