package expo

import (
	"testing"
)

func TestValidateResponseErrorStatus(t *testing.T) {
	response := &PushResponse{
		Status:  "error",
		Message: "failed",
		Details: map[string]string{},
	}
	err := response.ValidateResponse()
	typed, ok := err.(*PushResponseError)
	if !ok {
		t.Error("Incorrect error type")
	}
	if typed.Response != response {
		t.Error("Didn't return called response")
	}
}

func TestValidateResponseSuccess(t *testing.T) {
	response := &PushResponse{
		Status: "ok",
	}
	err := response.ValidateResponse()
	if err != nil {
		t.Error("Errored on valid response")
	}
}

func TestValidateResponseDeviceNotRegistered(t *testing.T) {
	response := &PushResponse{
		Status:  "error",
		Message: "Not registered",
		Details: map[string]string{"error": "DeviceNotRegistered"},
	}
	err := response.ValidateResponse()
	typed, ok := err.(*DeviceNotRegisteredError)
	if !ok {
		t.Error("Incorrect error type")
	}
	if typed.Response != response {
		t.Error("Didn't return called response")
	}
}

func TestValidateResponseErrorMessageTooBig(t *testing.T) {
	response := &PushResponse{
		Status:  "error",
		Message: "Message too big",
		Details: map[string]string{"error": "MessageTooBig"},
	}
	err := response.ValidateResponse()
	typed, ok := err.(*MessageTooBigError)
	if !ok {
		t.Error("Incorrect error type")
	}
	if typed.Response != response {
		t.Error("Didn't return called response")
	}
}

func TestValidateResponseErrorMessageRateExceeded(t *testing.T) {
	response := &PushResponse{
		Status:  "error",
		Message: "Too many messages at once",
		Details: map[string]string{"error": "MessageRateExceeded"},
	}
	err := response.ValidateResponse()
	typed, ok := err.(*MessageRateExceededError)
	if !ok {
		t.Error("Incorrect error type")
	}
	if typed.Response != response {
		t.Error("Didn't return called response")
	}
}
