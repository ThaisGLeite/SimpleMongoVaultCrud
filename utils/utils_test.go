package utils

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestHandleError(t *testing.T) {
	var buffer bytes.Buffer
	log.SetOutput(&buffer)

	HandleError("I", "Info Message", nil)
	if buffer.Len() != 0 {
		t.Errorf("Expected no log entry, got '%s'", buffer.String())
	}

	err := os.ErrNotExist
	HandleError("E", "Error Message", err)
	if !bytes.Contains(buffer.Bytes(), []byte("\"error\":\"file does not exist\"")) || !bytes.Contains(buffer.Bytes(), []byte("\"level\":\"error\"")) {
		t.Errorf("Expected error log entry, got '%s'", buffer.String())
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password1!", true},
		{"weak", false},
		{"NoSpecialChar1", false},
		{"Nouppercase1!", true},
		{"NOLOWERCASE1!", false},
		{"NoNumber!", false},
	}

	for _, test := range tests {
		actual := IsStrongPassword(test.password)
		if actual != test.expected {
			t.Errorf("For password '%s', expected '%v', got '%v'", test.password, test.expected, actual)
		} else {
			fmt.Printf("Password '%s' is strong: %v\n", test.password, actual)
		}
	}
}
