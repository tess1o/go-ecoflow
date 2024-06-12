package ecoflow

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHttpRequest_generateSign(t *testing.T) {
	tables := []struct {
		queryString  string
		nonce        string
		timestamp    string
		secretKey    string
		expectedSign string
	}{
		{
			queryString:  "",
			nonce:        "123456",
			timestamp:    "1671171709428",
			secretKey:    "secret",
			expectedSign: encryptHmacSHA256("accessKey=&nonce=123456&timestamp=1671171709428", "secret"),
		},
		{
			queryString:  "param1=v1&param2=v2",
			nonce:        "2345",
			timestamp:    "191817",
			secretKey:    "key",
			expectedSign: encryptHmacSHA256("param1=v1&param2=v2&accessKey=&nonce=2345&timestamp=191817", "key"),
		},
	}

	for _, table := range tables {
		req := &HttpRequest{
			secretKey:         table.secretKey,
			requestParameters: map[string]interface{}{},
		}
		actualSign := req.generateSign(table.queryString, table.nonce, table.timestamp)
		if actualSign != table.expectedSign {
			t.Errorf("GenerateSign incorrect. Got: %s, want: %s", actualSign, table.expectedSign)
		}
	}
}

func TestGetKeyValueString(t *testing.T) {
	tests := []struct {
		name          string
		queryString   string
		nonce         string
		timestamp     string
		accessKey     string
		secretKey     string
		expectedValue string
	}{
		{
			name:          "No Query String",
			queryString:   "",
			nonce:         "123",
			timestamp:     "1634796537",
			accessKey:     "ABCD",
			secretKey:     "",
			expectedValue: "accessKey=ABCD&nonce=123&timestamp=1634796537",
		},
		{
			name:          "With Query String",
			queryString:   "param1=value1&param2=value2",
			nonce:         "123",
			timestamp:     "1634796537",
			accessKey:     "ABCD",
			secretKey:     "",
			expectedValue: "param1=value1&param2=value2&accessKey=ABCD&nonce=123&timestamp=1634796537",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			request := &HttpRequest{
				accessKey: tt.accessKey,
				secretKey: tt.secretKey,
			}

			actualValue := request.getKeyValueString(tt.queryString, tt.nonce, tt.timestamp)
			if actualValue != tt.expectedValue {
				t.Fatalf("Expected: %s, but got: %s", tt.expectedValue, actualValue)
			}
		})
	}
}

func TestGenerateQueryParams(t *testing.T) {
	var tests = []struct {
		name        string
		input       map[string]interface{}
		expectedOut string
	}{
		{
			name:        "Simple Key Value Pairs",
			input:       map[string]interface{}{"key1": "value1", "key2": "value2"},
			expectedOut: "key1=value1&key2=value2",
		},
		{
			name: "Nested Key Value Pairs",
			input: map[string]interface{}{"key1": "value1", "key2": "value2", "abc": map[string]interface{}{
				"p1": "v1",
			}},
			expectedOut: "abc.p1=v1&key1=value1&key2=value2",
		},
		{
			name: "Complex value with nested arrays",
			input: map[string]interface{}{
				"name": "demo1",
				"ids":  []interface{}{1, 2, 3},
				"deviceInfo": map[string]interface{}{
					"id": 1,
				},
				"deviceList": []interface{}{
					map[string]interface{}{"id": 1},
					map[string]interface{}{"id": 2},
				},
			},
			expectedOut: "deviceInfo.id=1&deviceList[0].id=1&deviceList[1].id=2&ids[0]=1&ids[1]=2&ids[2]=3&name=demo1",
		},
		{
			name: "NestedData",
			input: map[string]interface{}{
				"commission": map[string]interface{}{
					"x": "y",
				},
			},
			expectedOut: "commission.x=y",
		},
		{
			name:        "With Numerical Values",
			input:       map[string]interface{}{"key1": 2, "key2": 10},
			expectedOut: "key1=2&key2=10",
		},
		{
			name:        "With Boolean Values",
			input:       map[string]interface{}{"key1": true, "key2": false},
			expectedOut: "key1=true&key2=false",
		},
		{
			name:        "With Empty Map",
			input:       map[string]interface{}{},
			expectedOut: "",
		},
		{
			name:        "With Nil Map",
			input:       nil,
			expectedOut: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := generateQueryParams(tt.input)
			if out != tt.expectedOut {
				t.Errorf("got %q, want %q", out, tt.expectedOut)
			}
		})
	}
}

func TestGenerateNonce(t *testing.T) {
	tests := []struct {
		name string
		want func(string) bool
	}{
		{
			name: "Should return a six digits string",
			want: func(got string) bool {
				return len(got) == 6
			},
		},
		{
			name: "Output should be integers",
			want: func(got string) bool {
				_, err := strconv.Atoi(got)
				return err == nil
			},
		},
		{
			name: "Consecutive calls should return different strings",
			want: func(got string) bool {
				return got != generateNonce()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateNonce(); !tt.want(got) {
				t.Errorf("generateNonce() = %v, want %v", got, tt.name)
			}
		})
	}
}

// Unit tests for the Execute method
func TestExecute(t *testing.T) {
	// Create a mock sign parameters function
	mockSignParams := signParameters{
		queryParams: "param1=value1&param2=value2",
		nonce:       "123456",
		timestamp:   fmt.Sprintf("%d", time.Now().Unix()),
		sign:        "signature",
	}
	getMockSignParams := func() *signParameters {
		return &mockSignParams
	}

	// Successful GET request
	t.Run("Successful GET request", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			// Validate headers
			if r.Header.Get(accessKeyHeader) != "testAccessKey" {
				t.Errorf("expected access key header to be 'testAccessKey', got '%s'", r.Header.Get(accessKeyHeader))
			}
			if r.Header.Get(nonceHeader) != mockSignParams.nonce {
				t.Errorf("expected nonce header to be '%s', got '%s'", mockSignParams.nonce, r.Header.Get(nonceHeader))
			}
			if r.Header.Get(timestampHeader) != mockSignParams.timestamp {
				t.Errorf("expected timestamp header to be '%s', got '%s'", mockSignParams.timestamp, r.Header.Get(timestampHeader))
			}
			if r.Header.Get(signHeader) != mockSignParams.sign {
				t.Errorf("expected sign header to be '%s', got '%s'", mockSignParams.sign, r.Header.Get(signHeader))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message": "success"}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		req := &HttpRequest{
			method:            http.MethodGet,
			uri:               server.URL,
			accessKey:         "testAccessKey",
			getSignParameters: getMockSignParams,
		}

		resp, err := req.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(resp) != `{"message": "success"}` {
			t.Errorf("expected response body to be '%s', got '%s'", `{"message": "success"}`, string(resp))
		}
	})

	// Successful POST request
	t.Run("Successful POST request", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			// Validate headers
			if r.Header.Get("Content-Type") != "application/json;charset=UTF-8" {
				t.Errorf("expected content-type header to be 'application/json;charset=UTF-8', got '%s'", r.Header.Get("Content-Type"))
			}
			if r.Header.Get(accessKeyHeader) != "testAccessKey" {
				t.Errorf("expected access key header to be 'testAccessKey', got '%s'", r.Header.Get(accessKeyHeader))
			}
			if r.Header.Get(nonceHeader) != mockSignParams.nonce {
				t.Errorf("expected nonce header to be '%s', got '%s'", mockSignParams.nonce, r.Header.Get(nonceHeader))
			}
			if r.Header.Get(timestampHeader) != mockSignParams.timestamp {
				t.Errorf("expected timestamp header to be '%s', got '%s'", mockSignParams.timestamp, r.Header.Get(timestampHeader))
			}
			if r.Header.Get(signHeader) != mockSignParams.sign {
				t.Errorf("expected sign header to be '%s', got '%s'", mockSignParams.sign, r.Header.Get(signHeader))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message": "success"}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		req := &HttpRequest{
			method:            http.MethodPost,
			uri:               server.URL,
			requestParameters: map[string]interface{}{"key": "value"},
			accessKey:         "testAccessKey",
			getSignParameters: getMockSignParams,
		}

		resp, err := req.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(resp) != `{"message": "success"}` {
			t.Errorf("expected response body to be '%s', got '%s'", `{"message": "success"}`, string(resp))
		}
	})

	// Invalid HTTP method
	t.Run("Invalid HTTP method", func(t *testing.T) {
		req := &HttpRequest{
			method:            "INVALID",
			uri:               "http://example.com",
			accessKey:         "testAccessKey",
			getSignParameters: getMockSignParams,
		}

		_, err := req.Execute(context.Background())
		if err == nil {
			t.Fatalf("expected error for invalid HTTP method, got nil")
		}
	})

	// HTTP request fails
	t.Run("HTTP request fails", func(t *testing.T) {
		req := &HttpRequest{
			method:            http.MethodGet,
			uri:               "http://invalid.url",
			accessKey:         "testAccessKey",
			getSignParameters: getMockSignParams,
		}

		_, err := req.Execute(context.Background())
		if err == nil {
			t.Fatalf("expected error for failed HTTP request, got nil")
		}
	})

	// HTTP response status is not OK
	t.Run("HTTP response status is not OK", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		req := &HttpRequest{
			method:            http.MethodGet,
			uri:               server.URL,
			accessKey:         "testAccessKey",
			getSignParameters: getMockSignParams,
		}

		_, err := req.Execute(context.Background())
		if err == nil {
			t.Fatalf("expected error for non-OK HTTP status, got nil")
		}
	})
}
