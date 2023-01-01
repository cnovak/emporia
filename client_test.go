package emporia

import (
	"fmt"
	"os"
	"testing"
)

func setupTestClient(t *testing.T) *Client {
	client, err := NewClient(os.Getenv("EMPORIA_USERNAME"), os.Getenv("EMPORIA_PASSWORD"))
	if err != nil {
		t.Fatalf(`NewClient() throws error: %v`, err)
	}
	return client
}

func ExampleNewClient() {
	username := os.Getenv("EMPORIA_USERNAME")
	password := os.Getenv("EMPORIA_PASSWORD")

	client, err := NewClient(username, password)
	if err != nil {
		panic(err)
	}

	devices, err := client.GetDevices()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Devices:\n%+v\n\n", devices)

}

func TestNewClient(t *testing.T) {
	setupTestClient(t)
}
