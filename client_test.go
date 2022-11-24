package emporia

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient(os.Getenv("EMPORIA_USERNAME"), os.Getenv("EMPORIA_PASSWORD"))
	if err != nil {
		t.Fatalf(`NewClient() throws error: %v`, err)
	}
}
