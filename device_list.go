package ecoflow

type DeviceListResponse struct {
	Code            string       `json:"code"`
	Message         string       `json:"message"`
	Devices         []DeviceInfo `json:"data"`
	EagleEyeTraceID string       `json:"eagleEyeTraceId"`
	Tid             string       `json:"tid"`
}

type DeviceInfo struct {
	SN     string `json:"sn"`
	Online int    `json:"online"`
}
