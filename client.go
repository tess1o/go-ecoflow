package ecoflow

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	ecoflowApiUrl  = "https://api.ecoflow.com"
	deviceListUrl  = ecoflowApiUrl + "/iot-open/sign/device/list"
	getAllQuoteUrl = ecoflowApiUrl + "/iot-open/sign/device/quota/all"
)

type Client struct {
	httpClient  *http.Client
	accessToken string
	secretToken string
}

func NewEcoflowClient(accessToken, secretToken string, httpClient *http.Client) *Client {
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

func (c *Client) GetDeviceList(ctx context.Context) (*DeviceListResponse, error) {
	request := NewHttpRequest(c.httpClient, "GET", deviceListUrl, nil, c.accessToken, c.secretToken)
	response, err := request.Execute(ctx)
	if err != nil {
		return nil, err
	}
	var deviceResponse DeviceListResponse

	err = json.Unmarshal(response, &deviceResponse)
	if err != nil {
		return nil, err
	}
	return &deviceResponse, nil
}

func (c *Client) GetDeviceAllQuote(ctx context.Context, deviceSn string) (*DeviceQuotaResponse, error) {
	params := make(map[string]interface{})
	params["sn"] = deviceSn

	request := NewHttpRequest(c.httpClient, "GET", getAllQuoteUrl, params, c.accessToken, c.secretToken)
	response, err := request.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var quotaResponse DeviceQuotaResponse

	err = json.Unmarshal(response, &quotaResponse)
	if err != nil {
		return nil, err
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
