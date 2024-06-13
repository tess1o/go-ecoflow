package ecoflow

import "testing"

func TestEncryptHmacSHA256(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		message  string
		secret   string
		expected string
	}{
		{
			name:     "Basic test",
			message:  "hello",
			secret:   "secret",
			expected: "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b", //calculated externally
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := encryptHmacSHA256(tc.message, tc.secret)
			if actual != tc.expected {
				t.Errorf("encryptHmacSHA256(%s, %s) = %s; expected %s", tc.message, tc.secret, actual, tc.expected)
			}
		})
	}
}
