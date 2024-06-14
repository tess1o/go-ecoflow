package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	ecoflowApiUrl  = "https://api.ecoflow.com"
	deviceListUrl  = ecoflowApiUrl + "/iot-open/sign/device/list"
	getAllQuoteUrl = ecoflowApiUrl + "/iot-open/sign/device/quota/all"
)

type Client struct {
	httpClient  *http.Client //can be customized if required
	accessToken string
	secretToken string
}

// NewEcoflowClient with default http client
func NewEcoflowClient(accessToken, secretToken string) *Client {
	return NewEcoflowClientWithHttpClient(accessToken, secretToken, nil)
}

// NewEcoflowClientWithHttpClient EcoflowClient with custom httpClient.
// HttpClient can be helpful when a proxy is required or some other custom transport, etc
func NewEcoflowClientWithHttpClient(accessToken, secretToken string, httpClient *http.Client) *Client {
	h := httpClient
	if h == nil {
		h = &http.Client{}
	}
	return &Client{
		httpClient:  h,
		accessToken: accessToken,
		secretToken: secretToken,
	}
}

// GetDeviceList executes a request to get the list of devises linked to the user account. Shared devices are not included
// If the response parameter "code" is not 0, then there is an error. Error code and error message are returned
func (c *Client) GetDeviceList(ctx context.Context) (*DeviceListResponse, error) {
	request := NewHttpRequest(c.httpClient, "GET", deviceListUrl, nil, c.accessToken, c.secretToken)
	response, err := request.Execute(ctx)
	if err != nil {
		return nil, err
	}
	var deviceResponse DeviceListResponse

	slog.Info("Response", string(response))

	err = json.Unmarshal(response, &deviceResponse)
	if err != nil {
		return nil, err
	}

	if deviceResponse.Code != "0" {
		return &deviceResponse, errors.New(fmt.Sprintf("can't get device list, error code: %s, error message: %s", deviceResponse.Code, deviceResponse.Message))
	}
	return &deviceResponse, nil
}

// GetDeviceAllQuote executes a request to get the parameters for a specific device.
// If the response parameter "code" is not 0, then there is an error. Error code and error message are returned.
// The raw parameters are unmarshalled into the DeviceQuotaResponse struct.
// The data field of the response is manually mapped to the appropriate structs (PdProperties, BmsEmsStatusProperties, BmsBmsStatusProperties, InvProperties, MpptProperties).
// The mapped data is assigned to the response struct and returned.
func (c *Client) GetDeviceAllQuote(ctx context.Context, deviceSn string) (*DeviceQuotaResponse, error) {
	response, err := c.getDeviceQuoteParams(ctx, deviceSn)
	if err != nil {
		return nil, err
	}

	var quotaResponse DeviceQuotaResponse

	err = json.Unmarshal(response, &quotaResponse)
	if err != nil {
		return nil, err
	}

	if quotaResponse.Code != "0" {
		return nil, errors.New(fmt.Sprintf("can't get parameters, error code: %s, error message: %s", quotaResponse.Code, quotaResponse.Message))
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(response, &jsonData)
	if err != nil {
		return nil, err
	}

	// Manually map the data field to the appropriate structs
	data := jsonData["data"].(map[string]interface{})

	// Convert data to JSON bytes to unmarshal into structs
	dataBytes, _ := json.Marshal(data)

	var pd PdProperties
	var bmsEmsStatus BmsEmsStatusProperties
	var bmsBmsStatus BmsBmsStatusProperties
	var inv InvProperties
	var mppt MpptProperties

	err = json.Unmarshal(dataBytes, &pd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &bmsEmsStatus)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &bmsBmsStatus)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &inv)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &mppt)
	if err != nil {
		return nil, err
	}

	// Assign to the response struct
	quotaResponse.Data.Pd = pd
	quotaResponse.Data.BmsEmsStatus = bmsEmsStatus
	quotaResponse.Data.BmsBmsStatus = bmsBmsStatus
	quotaResponse.Data.Inv = inv
	quotaResponse.Data.Mppt = mppt

	return &quotaResponse, nil
}

// GetDeviceQuoteRawParameters executes a request to get the raw parameters ("as is") for a specific device.
// It returns a map[string]interface{} containing the parameters and an error if any. The value type is mostly int, for some parameters it's float64 or []int
// If the response parameter "code" is not "0", then there is an error and the error message is returned.
// The parameters are taken from the Ecoflow response, "data" field
// Only parameters with float64 values are included in the returned map (int's are casted to float's)
// If the response is not valid or cannot be processed, an error is returned.
func (c *Client) GetDeviceQuoteRawParameters(ctx context.Context, deviceSn string) (map[string]interface{}, error) {
	response, err := c.getDeviceQuoteParams(ctx, deviceSn)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(response, &jsonData)
	if err != nil {
		return nil, err
	}

	if code, ok := jsonData["code"].(string); !ok || code != "0" {
		return nil, errors.New(fmt.Sprintf("can't get parameters, error code %s", code))
	}

	var params = make(map[string]interface{})
	dataMap, ok := jsonData["data"].(map[string]interface{})

	if !ok {
		return nil, errors.New("response is not valid, can't process it")
	}

	for k, v := range dataMap {
		if floatVal, isFloat := v.(float64); isFloat {
			params[k] = floatVal
		}
	}

	return params, err
}

// getDeviceQuoteParams executes a request to get the parameters for a specific device.
// The device serial number is passed as a parameter.
// The response is returned in the form of a byte array and an error, if any.
func (c *Client) getDeviceQuoteParams(ctx context.Context, deviceSn string) ([]byte, error) {
	params := make(map[string]interface{})
	params["sn"] = deviceSn

	request := NewHttpRequest(c.httpClient, "GET", getAllQuoteUrl, params, c.accessToken, c.secretToken)
	return request.Execute(ctx)
}
