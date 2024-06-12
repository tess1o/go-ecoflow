package ecoflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	accessKeyHeader = "accessKey"
	nonceHeader     = "nonce"
	timestampHeader = "timestamp"
	signHeader      = "sign"
)

type HttpRequest struct {
	httpClient        *http.Client
	method            string
	uri               string
	requestParameters map[string]interface{}
	accessKey         string
	secretKey         string
	getSignParameters func() *signParameters //required for unit testing
}

func NewHttpRequest(httpClient *http.Client, method string, uri string, params map[string]interface{}, accessKey, secretKey string) *HttpRequest {
	r := &HttpRequest{
		httpClient:        httpClient,
		method:            method,
		uri:               uri,
		requestParameters: params,
		accessKey:         accessKey,
		secretKey:         secretKey,
	}

	//required for unit testing
	r.getSignParameters = func() *signParameters {
		return r.generateSignParameters()
	}

	return r
}

func (r *HttpRequest) Execute(ctx context.Context) ([]byte, error) {
	signParams := r.getSignParameters()
	requestURI := r.uri + "?" + signParams.queryParams

	var reqBody bytes.Buffer

	if r.requestParameters != nil {
		reqBytes, _ := json.Marshal(r.requestParameters)
		reqBody.Write(reqBytes)
	}

	var httpReq *http.Request
	var err error

	switch r.method {
	case http.MethodGet:
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, requestURI, nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, r.uri, &reqBody)
		if err != nil {
			return nil, err
		}
		httpReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	default:
		slog.Error("Only POST and GET methods are supported so far")
		return nil, errors.New("unsupported http method")
	}

	httpReq.Header.Add(accessKeyHeader, r.accessKey)
	httpReq.Header.Add(nonceHeader, signParams.nonce)
	httpReq.Header.Add(timestampHeader, signParams.timestamp)
	httpReq.Header.Add(signHeader, signParams.sign)

	client := r.httpClient
	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response status is failed|url=%s, statusCode=%s", requestURI, resp.Status))
	}
	return io.ReadAll(resp.Body)
}

type signParameters struct {
	nonce       string
	timestamp   string
	accessKey   string
	sign        string
	queryParams string
}

func (r *HttpRequest) generateSignParameters() *signParameters {
	nonce := generateNonce()
	timestamp := generateTimestamp()
	queryParams := generateQueryParams(r.requestParameters)
	return &signParameters{
		queryParams: queryParams,
		nonce:       nonce,
		timestamp:   timestamp,
		accessKey:   r.accessKey,
		sign:        r.generateSign(queryParams, nonce, timestamp),
	}
}

// Step 4: encrypt
// E.g. byte[] signBytes = HMAC-SHA256(str, secretKey)
// Step 5: convert byte[] to hexadecimal string. String sign = bytesToHexString(signBytes)
// E.g. sign=85776ede686fe4783eac48135b0b1748ba2d7e9bb7791b826dc942fc29d4ada8
// Ecoflow documentation: https://developer-eu.ecoflow.com/us/document/generalInfo
func (r *HttpRequest) generateSign(queryString, nonce, timestamp string) string {
	keyValueString := r.getKeyValueString(queryString, nonce, timestamp)
	return encryptHmacSHA256(keyValueString, r.secretKey)
}

// The generate keyValue string that is used during generation of a "sing" header.
// The logic is to concatenate the values in specific order.
// From ecoflow documents
// Step 3: concatenate accessKey, nonce, timestamp
// E.g. str=param1=value1&param2=value2&accessKey=***&nonce=...&timestamp=...
// See step3 here: https://developer-eu.ecoflow.com/us/document/generalInfo
func (r *HttpRequest) getKeyValueString(queryString string, nonce string, timestamp string) string {
	keyValueString := accessKeyHeader + "=" + r.accessKey + "&" +
		nonceHeader + "=" + nonce + "&" +
		timestampHeader + "=" + timestamp

	if queryString != "" {
		keyValueString = queryString + "&" + keyValueString
	}
	return keyValueString
}

// From ecoflow documentation
// Step 1: request parameters must be sorted by ASCII value and concatenated with characters =, &
// E.g. str=param1=value1&param2=value2
// Step 2: if the type is nested, expand and splice according to the method of step 1.
// E.g. deviceInfo.id=1&deviceList[0].id=1&deviceList[1].id=2&ids[0]=1&ids[1]=2&ids[2]=3&name=demo1
// See step 1 and step 2 here: https://developer-eu.ecoflow.com/us/document/generalInfo
func generateQueryParams(data map[string]interface{}) string {
	var result []string

	// Process top-level map keys
	for k, v := range data {
		result = append(result, processValue(k, v)...)
	}

	// Sort results by ASCII value
	sort.Strings(result)

	// Concatenate results with & separator
	return strings.Join(result, "&")
}

func processValue(prefix string, value interface{}) []string {
	var result []string
	switch v := value.(type) {
	case map[string]interface{}:
		for k, nestedValue := range v {
			// Recursively process nested maps
			nestedPrefix := prefix + "." + k
			result = append(result, processValue(nestedPrefix, nestedValue)...)
		}
	case []interface{}:
		for i, item := range v {
			// Recursively process items in arrays
			nestedPrefix := prefix + "[" + strconv.Itoa(i) + "]"
			result = append(result, processValue(nestedPrefix, item)...)
		}
	case string:
		result = append(result, prefix+"="+v)
	case int:
		result = append(result, prefix+"="+strconv.Itoa(v))
	case float64:
		result = append(result, prefix+"="+strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		result = append(result, prefix+"="+strconv.FormatBool(v))
	}
	return result
}

// timestamp is a UTC timestamp (in nano)
func generateTimestamp() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// nonce is a random int with 6 digits
func generateNonce() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
