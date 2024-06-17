package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Ecoflow documentation:
// https://developer-eu.ecoflow.com/us/document/delta2
// https://developer-eu.ecoflow.com/us/document/delta2max

// TODO: Set AC discharge ("enabled" and X-Boost switch settings)
// TODO: { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"acOutCfg", "params":{ "enabled":0, "xboost":0, "out_voltage":30, "out_freq":1 } }

type PowerStation struct {
	c  *Client
	sn string
}

func (s *PowerStation) GetSn() string {
	return s.sn
}

type PowerStationModuleType int

const (
	PowerStationModuleTypePd       PowerStationModuleType = 1
	PowerStationModuleTypeBms      PowerStationModuleType = 2
	PowerStationModuleTypeInv      PowerStationModuleType = 3
	PowerStationModuleTypeBmsSlave PowerStationModuleType = 4
	PowerStationModuleTypeMppt     PowerStationModuleType = 5
)

//MPPT parameters

func (s *PowerStation) SetCarChargerSwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "mpptCar", PowerStationModuleTypeMppt, params)
}

func (s *PowerStation) SetBuzzerSilentMode(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "quietMode", PowerStationModuleTypeMppt, params)
}

// SetAcChargingSettings From ecoflow documentation:
// AC charging settings(chgPauseFlag: 0: AC charging in normal operation, 1: AC charging paused (not saved, restored by plugging))
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"acChgCfg", "params":{ "chgWatts":100, "chgPauseFlag":0 } }
// Most likely you want to use chgPauseFlag as 0 and set chargeWatts
func (s *PowerStation) SetAcChargingSettings(ctx context.Context, chargeWatts int, chgPauseFlag SettingSwitcher) (*CmdSetResponse, error) {
	if chargeWatts < 0 {
		return nil, errors.New("chargeWatts must be positive")
	}
	params := make(map[string]interface{})
	params["chgWatts"] = chargeWatts
	params["chgPauseFlag"] = chgPauseFlag
	return s.setParameter(ctx, "acChgCfg", PowerStationModuleTypeMppt, params)
}

// SetAcStandByTime AC standby time when there is no load(0: never shuts down, default value: 12 x 60 mins, unit: minute)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"standbyTime", "params":{ "standbyMins":180 } }
func (s *PowerStation) SetAcStandByTime(ctx context.Context, standbyMins int) (*CmdSetResponse, error) {
	if standbyMins < 0 {
		return nil, errors.New("standbyMins must be positive")
	}
	params := make(map[string]interface{})
	params["standbyMins"] = standbyMins
	return s.setParameter(ctx, "standbyTime", PowerStationModuleTypeMppt, params)
}

// SetCarStandByTime CAR standby duration settings(Auto shutdown when there is no load, unit: minute)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"carStandby", "params":{ "standbyMins":240 } }
func (s *PowerStation) SetCarStandByTime(ctx context.Context, standbyMins int) (*CmdSetResponse, error) {
	if standbyMins < 0 {
		return nil, errors.New("standbyMins must be positive")
	}
	params := make(map[string]interface{})
	params["standbyMins"] = standbyMins
	return s.setParameter(ctx, "carStandby", PowerStationModuleTypeMppt, params)
}

// Set12VDcChargingCurrent Set 12 V DC (car charger) charging current(Maximum DC charging current (mA), range: 4000 mA–10000 mA, default value: 8000 mA)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"dcChgCfg", "params":{ "dcChgCfg":5000 } }
func (s *PowerStation) Set12VDcChargingCurrent(ctx context.Context, chargingCurrent int) (*CmdSetResponse, error) {
	if chargingCurrent < 4000 || chargingCurrent > 100000 {
		return nil, errors.New("chargingCurrent out of range. Range: 4000 mA–10000 mA")
	}
	params := make(map[string]interface{})
	params["dcChgCfg"] = chargingCurrent
	return s.setParameter(ctx, "dcChgCfg", PowerStationModuleTypeMppt, params)
}

// PD parameters

// SetStandByTime Set standby time(0 for never standby; other values indicate the standby time; in minutes)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"standbyTime", "params":{ "standbyMin":0 } }
func (s *PowerStation) SetStandByTime(ctx context.Context, standbyMin int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["standbyMin"] = standbyMin
	return s.setParameter(ctx, "standbyTime", PowerStationModuleTypePd, params)
}

// SetDcSwitch Set DC(USB) switch(0: off, 1: on)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"dcOutCfg", "params":{ "enabled":0 } }
func (s *PowerStation) SetDcSwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "dcOutCfg", PowerStationModuleTypePd, params)
}

// SetLcdScreenTimeout LCD screen settings(delayOff: screen timeout, unit: seconds;brightLevel: must be set to 3; other values are invalid.)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"lcdCfg", "params":{ "delayOff":60, "brighLevel":3 } }
func (s *PowerStation) SetLcdScreenTimeout(ctx context.Context, delayOffSeconds int) (*CmdSetResponse, error) {
	if delayOffSeconds < 0 {
		return nil, errors.New("delayOff must be positive")
	}
	params := make(map[string]interface{})
	params["delayOff"] = delayOffSeconds
	params["brighLevel"] = 3
	return s.setParameter(ctx, "lcdCfg", PowerStationModuleTypePd, params)
}

// SetPrioritizePolarCharging Prioritize solar charging
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"pvChangePrio", "params":{ "pvChangeSet":0 } }
func (s *PowerStation) SetPrioritizePolarCharging(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "pvChangePrio", PowerStationModuleTypePd, params)
}

//TODO: Set energy management(isConfig: energy management, 0: disabled, 1: enabled; bpPowerSoc: backup reserve level; minDsgSoc: discharge limit (not in use);minChgSoc: charge limit (not in use))
//TODO: Set AC always on (acAutoOutConfig: 0: disabled; 1: enabled;minAcOutSoc: minimum SoC for turning on "AC always on" )

//BMS parameters

// SetMaxChargeSoC UPS settings(UPS, upper SoC limit when charging)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":2, "operateType":"upsConfig", "params":{ "maxChgSoc":50 } }
func (s *PowerStation) SetMaxChargeSoC(ctx context.Context, maxChgSoc int) (*CmdSetResponse, error) {
	if maxChgSoc < 0 {
		return nil, errors.New("maxChgSoc must be positive")
	}
	params := make(map[string]interface{})
	params["maxChgSoc"] = maxChgSoc
	return s.setParameter(ctx, "upsConfig", PowerStationModuleTypeBms, params)
}

// SetMinDischargeSoC SOC lower limit when discharging
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":2, "operateType":"dsgCfg", "params":{ "minDsgSoc":19 } }
func (s *PowerStation) SetMinDischargeSoC(ctx context.Context, minDsgSoc int) (*CmdSetResponse, error) {
	if minDsgSoc < 0 {
		return nil, errors.New("minDsgSoc must be positive")
	}
	params := make(map[string]interface{})
	params["minDsgSoc"] = minDsgSoc
	return s.setParameter(ctx, "dsgCfg", PowerStationModuleTypeBms, params)
}

// SetSoCToTurnOnSmartGenerator SoC that triggers EMS to turn on Smart Generator
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":2, "operateType":"openOilSoc", "params":{ "openOilSoc":40 } }
func (s *PowerStation) SetSoCToTurnOnSmartGenerator(ctx context.Context, openOilSoc int) (*CmdSetResponse, error) {
	if openOilSoc < 0 {
		return nil, errors.New("openOilSoc must be positive")
	}
	params := make(map[string]interface{})
	params["openOilSoc"] = openOilSoc
	return s.setParameter(ctx, "openOilSoc", PowerStationModuleTypeBms, params)
}

// SetSoCToTurnOffSmartGenerator SOC that triggers EMS to turn off Smart Generator
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":2, "operateType":"closeOilSoc", "params":{ "closeOilSoc":80 } }
func (s *PowerStation) SetSoCToTurnOffSmartGenerator(ctx context.Context, closeOilSoc int) (*CmdSetResponse, error) {
	if closeOilSoc < 0 {
		return nil, errors.New("closeOilSoc must be positive")
	}
	params := make(map[string]interface{})
	params["closeOilSoc"] = closeOilSoc
	return s.setParameter(ctx, "closeOilSoc", PowerStationModuleTypeBms, params)
}

func (s *PowerStation) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

func (s *PowerStation) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

func (s *PowerStation) setParameter(ctx context.Context, opType string, modType PowerStationModuleType, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:          fmt.Sprint(time.Now().UnixMilli()),
		OperateType: opType,
		ModuleType:  modType,
		Sn:          s.sn,
		Params:      params,
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
