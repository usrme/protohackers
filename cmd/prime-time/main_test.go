package main

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"prime number 2", 2, true},
		{"prime number 3", 3, true},
		{"prime number 5", 5, true},
		{"prime number 7", 7, true},
		{"prime number 11", 11, true},
		{"prime number 97", 97, true},
		{"composite number 4", 4, false},
		{"composite number 6", 6, false},
		{"composite number 8", 8, false},
		{"composite number 9", 9, false},
		{"composite number 15", 15, false},
		{"number 1", 1, false},
		{"number 0", 0, false},
		{"negative number -5", -5, false},
		{"negative number -2", -2, false},
		{"floating point 2.5", 2.5, false},
		{"floating point 3.14", 3.14, false},
		{"floating point 2.0", 2.0, true},
		{"floating point 3.0", 3.0, true},
		{"large prime 982451653", 982451653, true},
		{"large composite 982451654", 982451654, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPrime(tt.input)
			if result != tt.expected {
				t.Errorf("isPrime(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsMalformed(t *testing.T) {
	tests := []struct {
		name     string
		req      request
		expected bool
	}{
		{
			name:     "valid request",
			req:      request{Method: "isPrime", Number: floatPtr(5)},
			expected: false,
		},
		{
			name:     "nil number",
			req:      request{Method: "isPrime", Number: nil},
			expected: true,
		},
		{
			name:     "wrong method",
			req:      request{Method: "notPrime", Number: floatPtr(5)},
			expected: true,
		},
		{
			name:     "empty method",
			req:      request{Method: "", Number: floatPtr(5)},
			expected: true,
		},
		{
			name:     "wrong method with nil number",
			req:      request{Method: "wrongMethod", Number: nil},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMalformed(tt.req)
			if result != tt.expected {
				t.Errorf("isMalformed(%+v) = %v, expected %v", tt.req, result, tt.expected)
			}
		})
	}
}

func TestTCPHandler(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handler(conn)
		}
	}()

	addr := listener.Addr().String()

	tests := []struct {
		name           string
		input          string
		expectedOutput string
		shouldClose    bool
	}{
		{
			name:           "valid prime request",
			input:          `{"method":"isPrime","number":7}`,
			expectedOutput: `{"method":"isPrime","prime":true}`,
			shouldClose:    false,
		},
		{
			name:           "valid composite request",
			input:          `{"method":"isPrime","number":8}`,
			expectedOutput: `{"method":"isPrime","prime":false}`,
			shouldClose:    false,
		},
		{
			name:           "valid floating point request",
			input:          `{"method":"isPrime","number":2.5}`,
			expectedOutput: `{"method":"isPrime","prime":false}`,
			shouldClose:    false,
		},
		{
			name:           "valid integer as float request",
			input:          `{"method":"isPrime","number":2.0}`,
			expectedOutput: `{"method":"isPrime","prime":true}`,
			shouldClose:    false,
		},
		{
			name:        "malformed JSON",
			input:       `{"method":"isPrime","number":}`,
			shouldClose: true,
		},
		{
			name:        "missing number field",
			input:       `{"method":"isPrime"}`,
			shouldClose: true,
		},
		{
			name:        "wrong method",
			input:       `{"method":"notPrime","number":5}`,
			shouldClose: true,
		},
		{
			name:        "missing method field",
			input:       `{"number":5}`,
			shouldClose: true,
		},
		{
			name:        "number field is string",
			input:       `{"method":"isPrime","number":"5"}`,
			shouldClose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Fatalf("Failed to connect: %v", err)
			}
			defer conn.Close()

			// Send request
			_, err = conn.Write([]byte(tt.input + "\n"))
			if err != nil {
				t.Fatalf("Failed to write: %v", err)
			}

			if tt.shouldClose {
				// For malformed requests, expect connection to close
				time.Sleep(100 * time.Millisecond)

				// Try to read, should get error or "malformed request"
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				scanner := bufio.NewScanner(conn)
				if scanner.Scan() {
					response := scanner.Text()
					if !strings.Contains(response, "malformed") {
						t.Errorf("Expected malformed response, got: %s", response)
					}
				}
				return
			}

			// Read response
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			scanner := bufio.NewScanner(conn)
			if !scanner.Scan() {
				t.Fatalf("Failed to read response: %v", scanner.Err())
			}

			response := scanner.Text()

			// Parse and compare JSON responses
			var expectedResp, actualResp map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expectedOutput), &expectedResp); err != nil {
				t.Fatalf("Failed to parse expected response: %v", err)
			}
			if err := json.Unmarshal([]byte(response), &actualResp); err != nil {
				t.Fatalf("Failed to parse actual response: %v", err)
			}

			if actualResp["method"] != expectedResp["method"] || actualResp["prime"] != expectedResp["prime"] {
				t.Errorf("Expected %+v, got %+v", expectedResp, actualResp)
			}
		})
	}
}

func TestMultipleRequests(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handler(conn)
		}
	}()

	addr := listener.Addr().String()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	requests := []struct {
		input    string
		expected string
	}{
		{`{"method":"isPrime","number":2}`, `{"method":"isPrime","prime":true}`},
		{`{"method":"isPrime","number":4}`, `{"method":"isPrime","prime":false}`},
		{`{"method":"isPrime","number":17}`, `{"method":"isPrime","prime":true}`},
	}

	for i, req := range requests {
		// Send request
		_, err = conn.Write([]byte(req.input + "\n"))
		if err != nil {
			t.Fatalf("Request %d: Failed to write: %v", i, err)
		}

		// Read response
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		scanner := bufio.NewScanner(conn)
		if !scanner.Scan() {
			t.Fatalf("Request %d: Failed to read response: %v", i, scanner.Err())
		}

		response := scanner.Text()

		// Parse and compare JSON responses
		var expectedResp, actualResp map[string]interface{}
		if err := json.Unmarshal([]byte(req.expected), &expectedResp); err != nil {
			t.Fatalf("Request %d: Failed to parse expected response: %v", i, err)
		}
		if err := json.Unmarshal([]byte(response), &actualResp); err != nil {
			t.Fatalf("Request %d: Failed to parse actual response: %v", i, err)
		}

		if actualResp["method"] != expectedResp["method"] || actualResp["prime"] != expectedResp["prime"] {
			t.Errorf("Request %d: Expected %+v, got %+v", i, expectedResp, actualResp)
		}
	}
}

func TestConcurrentConnections(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handler(conn)
		}
	}()

	addr := listener.Addr().String()

	// Test 5 concurrent connections
	numConnections := 5
	done := make(chan bool, numConnections)

	for i := 0; i < numConnections; i++ {
		go func(clientID int) {
			defer func() { done <- true }()

			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Errorf("Client %d: Failed to connect: %v", clientID, err)
				return
			}
			defer conn.Close()

			request := `{"method":"isPrime","number":7}`
			expected := `{"method":"isPrime","prime":true}`

			_, err = conn.Write([]byte(request + "\n"))
			if err != nil {
				t.Errorf("Client %d: Failed to write: %v", clientID, err)
				return
			}

			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			scanner := bufio.NewScanner(conn)
			if !scanner.Scan() {
				t.Errorf("Client %d: Failed to read response: %v", clientID, scanner.Err())
				return
			}

			response := scanner.Text()

			var expectedResp, actualResp map[string]interface{}
			if err := json.Unmarshal([]byte(expected), &expectedResp); err != nil {
				t.Errorf("Client %d: Failed to parse expected response: %v", clientID, err)
				return
			}
			if err := json.Unmarshal([]byte(response), &actualResp); err != nil {
				t.Errorf("Client %d: Failed to parse actual response: %v", clientID, err)
				return
			}

			if actualResp["method"] != expectedResp["method"] || actualResp["prime"] != expectedResp["prime"] {
				t.Errorf("Client %d: Expected %+v, got %+v", clientID, expectedResp, actualResp)
			}
		}(i)
	}

	// Wait for all clients to complete
	for i := 0; i < numConnections; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent connections to complete")
		}
	}
}

func TestExtraneousFields(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handler(conn)
		}
	}()

	addr := listener.Addr().String()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Request with extraneous fields should be accepted
	request := `{"method":"isPrime","number":7,"extra":"ignored","another":123}`
	expected := `{"method":"isPrime","prime":true}`

	_, err = conn.Write([]byte(request + "\n"))
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		t.Fatalf("Failed to read response: %v", scanner.Err())
	}

	response := scanner.Text()

	var expectedResp, actualResp map[string]interface{}
	if err := json.Unmarshal([]byte(expected), &expectedResp); err != nil {
		t.Fatalf("Failed to parse expected response: %v", err)
	}
	if err := json.Unmarshal([]byte(response), &actualResp); err != nil {
		t.Fatalf("Failed to parse actual response: %v", err)
	}

	if actualResp["method"] != expectedResp["method"] || actualResp["prime"] != expectedResp["prime"] {
		t.Errorf("Expected %+v, got %+v", expectedResp, actualResp)
	}
}

// Helper function to create float64 pointer
func floatPtr(f float64) *float64 {
	return &f
}
