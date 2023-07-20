package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	parser := NewParser()

	// Test 1: Valid JSON-RPC message
	message, err := parser.Parse(strings.NewReader("Content-Length: 42\r\n\r\n{\"jsonrpc\": \"2.0\", \"method\": \"initialize\"}t"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := Request {
    Jsonrpc: "2.0",
    Method:  "initialize",
  }

	if message.Jsonrpc != expected.Jsonrpc || message.Method != expected.Method || message.ID != expected.ID {
		t.Errorf("Expected '%v', got '%v'", expected, message)
	}

	// Test 2: Invalid header
	_, err = parser.Parse(strings.NewReader("Invalid-Header: 43\r\n\r\n{\"jsonrpc\": \"2.0\", \"method\": \"initialize\"}"))
	if err == nil {
		t.Error("Expected error for invalid header, got nil")
	}

	// Test 3: No header
	_, err = parser.Parse(strings.NewReader("{\"jsonrpc\": \"2.0\", \"method\": \"initialize\"}"))
	if err == nil {
		t.Error("Expected error for missing header, got nil")
	}

	// Test 4: No body
	_, err = parser.Parse(strings.NewReader("Content-Length: 43\r\n\r\n"))
	if err == nil {
		t.Error("Expected error for missing body, got nil")
	}

	// Test 5: Invalid Content-Length
	_, err = parser.Parse(strings.NewReader("Content-Length: invalid\r\n\r\n{\"jsonrpc\": \"2.0\", \"method\": \"initialize\"}"))
	if err == nil {
		t.Error("Expected error for invalid content length, got nil")
	}

	// Test 6: Content-Length does not match body length
	_, err = parser.Parse(strings.NewReader("Content-Length: 44\r\n\r\n{\"jsonrpc\": \"2.0\", \"method\": \"initialize\"}"))
	if err == nil {
		t.Error("Expected error for mismatched content length, got nil")
	}
}

