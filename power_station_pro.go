package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Ecoflow documentation: https://developer-eu.ecoflow.com/us/document/deltapro
// The API for "regular" power stations (like Delta 2, Delta 2 Max, River 2, etc) is different from the "PRO" version

type PowerStationPro struct {
	c  *Client
	sn string
}

func (s *PowerStationPro) GetSn() string {
	return s.sn
}

// SetXboostSwitcher Setting the X-Boost switch
// "params":{ "cmdSet": 32, "id": 66, "enabled": 0, "xboost": 0 }
func (s *PowerStationPro) SetXboostSwitcher(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 66
	params["enabled"] = enabled
	params["xboost"] = enabled
	return s.setParameter(ctx, params)
}

// SetCarChargerSwitch Setting the car charger switch
// "params":{ "cmdSet": 32, "id": 81, "enabled": 1 }
func (s *PowerStationPro) SetCarChargerSwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 81
	params["enabled"] = enabled
	return s.setParameter(ctx, params)
}

// SetMaxChargeLevel Setting the charge level
// "params":{ "cmdSet": 32, "id": 49, "maxChgSoc": 100 }
func (s *PowerStationPro) SetMaxChargeLevel(ctx context.Context, maxChgSoc int) (*CmdSetResponse, error) {
	if maxChgSoc < 0 || maxChgSoc > 100 {
		return nil, errors.New("maxChgSoc out of range. Range 0:100")
	}
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 49
	params["maxChgSoc"] = maxChgSoc
	return s.setParameter(ctx, params)
}

// SetMinDischargeLevel Setting the discharge level
// "params":{ "cmdSet": 32, "id": 51, "minDsgSoc": 10 }
func (s *PowerStationPro) SetMinDischargeLevel(ctx context.Context, minDsgSoc int) (*CmdSetResponse, error) {
	if minDsgSoc < 0 || minDsgSoc > 100 {
		return nil, errors.New("minDsgSoc out of range. Range 0:100")
	}
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 51
	params["minDsgSoc"] = minDsgSoc
	return s.setParameter(ctx, params)
}

// SetCarInputCurrent Setting the car input current
// "params":{ "cmdSet": 32, "id": 71, "currMa": 4000 }
func (s *PowerStationPro) SetCarInputCurrent(ctx context.Context, currMa int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 71
	params["currMa"] = currMa
	return s.setParameter(ctx, params)
}

// SetBeepSwitch Setting the beep switch
// "params":{ "cmdSet": 32, "id": 38, "enabled": 1 }
func (s *PowerStationPro) SetBeepSwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 38
	params["enabled"] = enabled
	return s.setParameter(ctx, params)
}

// SetScreenBrightness Setting the screen brightness
// "params":{ "cmdSet": 32, "id": 39, "lcdBrightness": 100 }
func (s *PowerStationPro) SetScreenBrightness(ctx context.Context, lcdBrightness int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 39
	params["lcdBrightness"] = lcdBrightness
	return s.setParameter(ctx, params)
}

// SetSoCToTurnOnSmartGenerator Setting the lower threshold percentage of smart generator auto on
// "params":{ "cmdSet": 32, "id": 52, "openOilSoc": 52 }
func (s *PowerStationPro) SetSoCToTurnOnSmartGenerator(ctx context.Context, openOilSoc int) (*CmdSetResponse, error) {
	if openOilSoc < 0 {
		return nil, errors.New("openOilSoc must be positive")
	}
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 52
	params["openOilSoc"] = openOilSoc
	return s.setParameter(ctx, params)
}

// SetSoCToTurnOffSmartGenerator Setting the upper threshold percentage of smart generator auto off
// "params":{ "cmdSet": 32, "id": 53, "closeOilSoc": 10 }
func (s *PowerStationPro) SetSoCToTurnOffSmartGenerator(ctx context.Context, closeOilSoc int) (*CmdSetResponse, error) {
	if closeOilSoc < 0 {
		return nil, errors.New("closeOilSoc must be positive")
	}
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 53
	params["closeOilSoc"] = closeOilSoc
	return s.setParameter(ctx, params)
}

// SetUnitTimeout Setting the unit timeout
// "params":{ "cmdSet": 32, "id": 33, "standByMode": 0 }
func (s *PowerStationPro) SetUnitTimeout(ctx context.Context, standByMode int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 33
	params["standByMode"] = standByMode
	return s.setParameter(ctx, params)
}

// SetScreenTimeout Setting the screen timeout
// "params":{ "cmdSet": 32, "id": 39, "lcdTime": 60 }
func (s *PowerStationPro) SetScreenTimeout(ctx context.Context, lcdTime int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 39
	params["lcdTime"] = lcdTime
	return s.setParameter(ctx, params)
}

// SetAcStandByTime Setting the AC standby time
// "params":{ "cmdSet": 32, "id": 153, "standByMins": 720 }
func (s *PowerStationPro) SetAcStandByTime(ctx context.Context, standByMins int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 153
	params["standByMins"] = standByMins
	return s.setParameter(ctx, params)
}

// SetAcChargingSettings AC charging settings
// "params":{ "cmdSet": 32, "id": 69, "slowChgPower": 0 }
func (s *PowerStationPro) SetAcChargingSettings(ctx context.Context, slowChgPower int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 69
	params["slowChgPower"] = slowChgPower
	return s.setParameter(ctx, params)
}

// SetPvChargingType PV charging type
// "params":{ "cmdSet": 32, "id": 82, "chgType": 0 }
func (s *PowerStationPro) SetPvChargingType(ctx context.Context, chgType PowerStationPvChargeType) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 82
	params["chgType"] = chgType
	return s.setParameter(ctx, params)
}

// SetBypassAcAutoStart Bypass AC auto start
// "params":{ "cmdSet": 32, "id": 84, "enabled": 0 }
func (s *PowerStationPro) SetBypassAcAutoStart(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["cmdSet"] = 32
	params["id"] = 84
	params["enabled"] = enabled
	return s.setParameter(ctx, params)
}

func (s *PowerStationPro) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

func (s *PowerStationPro) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

func (s *PowerStationPro) setParameter(ctx context.Context, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:     fmt.Sprint(time.Now().UnixMilli()),
		Sn:     s.sn,
		Params: params,
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
