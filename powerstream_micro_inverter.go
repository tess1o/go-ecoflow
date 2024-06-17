package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Ecoflow documentation:
// https://developer-eu.ecoflow.com/us/document/powerStreamMicroInverter

type PowerStreamMicroInverter struct {
	c  *Client
	sn string
}

func (s *PowerStreamMicroInverter) GetSn() string {
	return s.sn
}

// SetPowerSupplyPriority Power supply priority settings(0: prioritize power supply; 1: prioritize power storage)
// {"sn": "HW513000SF767194","cmdCode": "WN511_SET_SUPPLY_PRIORITY_PACK","params": {"supplyPriority": 0}}
func (s *PowerStreamMicroInverter) SetPowerSupplyPriority(ctx context.Context, supplyPriority int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["supplyPriority"] = supplyPriority
	return s.setParameter(ctx, "WN511_SET_SUPPLY_PRIORITY_PACK", params)
}

// SetCustomLoadPowerSettings Custom load power settings(Range: 0 W–600 W; unit: 0.1 W)
// {"sn": "HW513000SF767194","cmdCode": "WN511_SET_PERMANENT_WATTS_PACK","params": {"permanentWatts": 20}}
func (s *PowerStreamMicroInverter) SetCustomLoadPowerSettings(ctx context.Context, permanentWatts float64) (*CmdSetResponse, error) {
	if permanentWatts < 0 || permanentWatts > 600 {
		return nil, errors.New("permanentWatts is out of range. Range 0:600, unit 0.1W")
	}
	params := make(map[string]interface{})
	params["permanentWatts"] = permanentWatts
	return s.setParameter(ctx, "WN511_SET_PERMANENT_WATTS_PACK", params)
}

// SetLowerLimitSettingsForBatterDischarging Lower limit settings for battery discharging(lowerLimit: 1-30)
// {"sn": "HW513000SF767194","cmdCode": "WN511_SET_BAT_LOWER_PACK","params": {"lowerLimit": 20}}
func (s *PowerStreamMicroInverter) SetLowerLimitSettingsForBatterDischarging(ctx context.Context, lowerLimit float64) (*CmdSetResponse, error) {
	if lowerLimit < 1 || lowerLimit > 30 {
		return nil, errors.New("lowerLimit is out of range. Range 1:30")
	}
	params := make(map[string]interface{})
	params["lowerLimit"] = lowerLimit
	return s.setParameter(ctx, "WN511_SET_BAT_LOWER_PACK", params)
}

// SetUpperLimitSettingsForBatterCharging Upper limit settings for battery charging(upperLimit: 70-100)
// {"sn": "HW513000SF767194","cmdCode": "WN511_SET_BAT_UPPER_PACK","params": {"upperLimit": 80}}
func (s *PowerStreamMicroInverter) SetUpperLimitSettingsForBatterCharging(ctx context.Context, upperLimit float64) (*CmdSetResponse, error) {
	if upperLimit < 70 || upperLimit > 100 {
		return nil, errors.New("upperLimit is out of range. Range 70:100")
	}
	params := make(map[string]interface{})
	params["upperLimit"] = upperLimit
	return s.setParameter(ctx, "WN511_SET_BAT_UPPER_PACK", params)
}

// SetLightBrightness Indicator light brightness adjustment(rgb brightness: 0-1023 (the larger the value, the higher the brightness); default value: 1023)
// {"sn": "HW513000SF767194","cmdCode": "WN511_SET_BRIGHTNESS_PACK","params": {"brightness": 200}}
func (s *PowerStreamMicroInverter) SetLightBrightness(ctx context.Context, brightness float64) (*CmdSetResponse, error) {
	if brightness < 0 || brightness > 1023 {
		return nil, errors.New("brightness is out of range. Range 0:1023")
	}
	params := make(map[string]interface{})
	params["brightness"] = brightness
	return s.setParameter(ctx, "WN511_SET_BRIGHTNESS_PACK", params)
}

// DeleteScheduledSwitchingTasks Deleting scheduled switching tasks(taskIndex: 0-10)
// {"sn": "HW513000SF767194","cmdCode": "WN511_DELETE_TIME_TASK","params": {"taskIndex": 1}}
func (s *PowerStreamMicroInverter) DeleteScheduledSwitchingTasks(ctx context.Context, taskIndex float64) (*CmdSetResponse, error) {
	if taskIndex < 0 || taskIndex > 10 {
		return nil, errors.New("taskIndex is out of range. Range 0:10")
	}
	params := make(map[string]interface{})
	params["taskIndex"] = taskIndex
	return s.setParameter(ctx, "WN511_DELETE_TIME_TASK", params)
}

func (s *PowerStreamMicroInverter) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

func (s *PowerStreamMicroInverter) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

func (s *PowerStreamMicroInverter) setParameter(ctx context.Context, cmdCode string, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:      fmt.Sprint(time.Now().UnixMilli()),
		CmdCode: cmdCode,
		Sn:      s.sn,
		Params:  params,
	}

	jsonData, err := json.Marshal(cmdReq)
	if err != nil {
		return nil, err
	}

	var req map[string]interface{}

	err = json.Unmarshal(jsonData, &req)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, req)
}
