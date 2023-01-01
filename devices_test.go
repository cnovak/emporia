package emporia

import (
	"testing"
)

func TestClient_GetDevices(t *testing.T) {

	c := setupTestClient(t)

	got, err := c.GetDevices()
	if err != nil {
		t.Errorf("Client.GetDevices() error = %v", err)
		return
	}

	if got == nil {
		t.Errorf("Client.GetDevices() got = %v", got)
		return
	}

}

func TestClient_GetSolarDevices(t *testing.T) {

	c := setupTestClient(t)

	got, err := c.GetSolar()
	if err != nil {
		t.Errorf("Client.GetSolar() error = %v", err)
		return
	}

	if got == nil {
		t.Error("Client.GetSolar() returned nil")
		return
	}

}
