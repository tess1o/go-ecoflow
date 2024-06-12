package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Method            string
	URI               string
	RequestParameters map[string]interface{}
	accessKey         string
	secretKey         string
}

func NewHttpRequest(method string, uri string, params map[string]interface{}, accessKey, secretKey string) *HttpRequest {
	return &HttpRequest{
		Method:            method,
		URI:               uri,
		RequestParameters: params,
		accessKey:         accessKey,
		secretKey:         secretKey,
	}
}

func (r *HttpRequest) Execute() ([]byte, error) {
	signParams := r.getSignParameters()
	requestURI := r.URI + "?" + signParams.queryParams

	var reqBody bytes.Buffer

	if r.RequestParameters != nil {
		reqBytes, _ := json.Marshal(r.RequestParameters)
		reqBody.Write(reqBytes)
	}

	var httpReq *http.Request
	var err error

	switch r.Method {
	case http.MethodGet:
		httpReq, err = http.NewRequest(http.MethodGet, requestURI, nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		httpReq, err = http.NewRequest(http.MethodPost, r.URI, &reqBody)
		if err != nil {
			return nil, err
		}
		httpReq.Header.Add("Content-Type", "application/json;charset=UTF-8")
	}

	httpReq.Header.Add(accessKeyHeader, r.accessKey)
	httpReq.Header.Add(nonceHeader, signParams.nonce)
	httpReq.Header.Add(timestampHeader, signParams.timestamp)
	httpReq.Header.Add(signHeader, signParams.sign)

	client := http.Client{}

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

func (r *HttpRequest) getSignParameters() *signParameters {
	nonce := generateNonce()
	timestamp := generateTimestamp()
	queryParams := generateQueryParams(r.RequestParameters)
	return &signParameters{
		queryParams: queryParams,
		nonce:       nonce,
		timestamp:   timestamp,
		accessKey:   r.accessKey,
		sign:        r.generateSign(queryParams, nonce, timestamp),
	}
}

func (r *HttpRequest) generateSign(queryString, nonce, timestamp string) string {
	keyValueString := r.getKeyValueString(queryString, nonce, timestamp)

	return encryptHmacSHA256(keyValueString, r.secretKey)
}

func (r *HttpRequest) getKeyValueString(queryString string, nonce string, timestamp string) string {
	keyValueString := accessKeyHeader + "=" + r.accessKey + "&" +
		nonceHeader + "=" + nonce + "&" +
		timestampHeader + "=" + timestamp

	if queryString != "" {
		keyValueString = queryString + "&" + keyValueString
	}
	return keyValueString
}

func generateQueryParams(requestParams map[string]interface{}) string {
	sortKeyValueMap := requestParams
	keys := make([]string, 0, len(sortKeyValueMap))
	for k := range sortKeyValueMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	queryString := ""
	for _, k := range keys {
		queryString += k + "=" + fmt.Sprint(sortKeyValueMap[k]) + "&"
	}
	queryString = strings.TrimRight(queryString, "&")
	return queryString
}

func generateTimestamp() string {
	return fmt.Sprint(time.Now().UnixNano())
}

func generateNonce() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
