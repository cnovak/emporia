package emporia

import (
	"testing"
	"time"
)

func TestClient_GetDeviceUsages(t *testing.T) {

	c := setupTestClient(t)

	devices, err := c.GetDevices()
	if err != nil {
		t.Errorf("Client.GetDevices() error = %v", err)
		return
	}

	if devices == nil {
		t.Errorf("Client.GetDevices() got = %v", devices)
		return
	}

	// Get meter device
	var energyMeterDevices []uint64
	for _, d := range devices.Devices {
		if d.Model == "VUE002" {
			energyMeterDevices = append(energyMeterDevices, d.DeviceGID)
		}
	}

	got, err := c.GetDeviceListUsages(energyMeterDevices, Second, KilowattHours, time.Now().UTC())
	if err != nil {
		t.Errorf("Client.GetDevices() error = %v", err)
		return
	}

	if got == nil {
		t.Errorf("Client.GetDevices() got = %v", got)
		return
	}

	if len(got.Devices) == 0 {
		t.Error("Client.GetDevices() is an empty list")
		return
	}

	// Get Net Usage channel
	NetUsageChannel, err := got.Devices[0].getChannel(NetUsage)
	if err != nil {
		t.Errorf("device.getChannel() error = %v", err)
		return
	}

	if NetUsageChannel == nil {
		t.Errorf("device.getChannel() is nil")
		return
	}

	// Get Net Usage channel
	TotalUsageChannel, err := got.Devices[0].getChannel(TotalUsage)
	if err != nil {
		t.Errorf("device.getChannel() error = %v", err)
		return
	}

	if TotalUsageChannel == nil {
		t.Errorf("device.getChannel() is nil")
		return
	}

}
