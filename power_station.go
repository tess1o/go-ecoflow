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
// For PRO version the API is different, probably a separate struct will be created

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

type PowerStationPvChargeType int

const (
	PowerStationPvChargeTypeAuto    PowerStationPvChargeType = 0
	PowerStationPvChargeTypeMppt    PowerStationPvChargeType = 1
	PowerStationPvChargeTypeAdapter PowerStationPvChargeType = 2
)

//MPPT parameters

func (s *PowerStation) SetBuzzerSilentMode(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "quietMode", PowerStationModuleTypeMppt, params)
}

func (s *PowerStation) SetCarChargerSwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := getParamsEnabled(enabled)
	return s.setParameter(ctx, "mpptCar", PowerStationModuleTypeMppt, params)
}

// SetAcEnabled Set AC discharge ("enabled" and X-Boost switch settings)
// AC discharging settings(enabled: AC switch, 0: off, 1: on; xboost: X-Boost switch, 0: off, 1: on; out_voltage: output voltage, read-only; out_freq: output frequency, 1: 50 Hz, 2: 60 Hz, other values are invalid)
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":5, "operateType":"acOutCfg", "params":{ "enabled":0, "xboost":0, "out_voltage":30, "out_freq":1 } }
// outFreq: 1 for 50Hz, 2 for 60Hz (check your grid), outVoltage 220 for Europe, for the USA probably 120
// It appears that all 4 parameters must be sent, otherwise it doesn't apply the changes
func (s *PowerStation) SetAcEnabled(ctx context.Context, acEnabled, xBoostEnabled SettingSwitcher, outFreq GridFrequency, outVoltage int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["enabled"] = acEnabled
	params["xboost"] = xBoostEnabled
	params["out_freq"] = outFreq
	params["out_voltage"] = outVoltage
	return s.setParameter(ctx, "acOutCfg", PowerStationModuleTypeMppt, params)
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

// SetPvChargingTypeSettings PV charging type settings(chaType: 0: auto identification, 1: MPPT, 2: adapter, other values are invalid; chaType2: 0: auto identification, 1: MPPT, 2: adapter, other values are invalid)
// {"id": 123,"version": "1.0","sn": "R351ZFB4HF6L0030","moduleType": 5,"operateType": "chaType","params": {"chaType": 0,"chaType2": 0}}
func (s *PowerStation) SetPvChargingTypeSettings(ctx context.Context, chargeType1, chargeType2 PowerStationPvChargeType) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["chaType"] = chargeType1
	params["chaType2"] = chargeType2
	return s.setParameter(ctx, "chaType", PowerStationModuleTypeMppt, params)
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

// SetEnergyManagement Set energy management(isConfig: energy management, 0: disabled, 1: enabled; bpPowerSoc: backup reserve level; minDsgSoc: discharge limit (not in use);minChgSoc: charge limit (not in use))
// Energy statistics(isConfig: energy management, 0: enabled, 1: disabled; bpPowerSoc: backup reserve level; minDsgSoc: lower limit when discharging (not in use); minChgSoc: upper limit when charging (not in use))
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"watthConfig", "params":{ "isConfig":1, "bpPowerSoc":95, "minDsgSoc":40, "minChgSoc":95 } }
func (s *PowerStation) SetEnergyManagement(ctx context.Context, enabled SettingSwitcher, bpPowerSoc, minDsgSoc, minChgSoc int) (*CmdSetResponse, error) {
	if bpPowerSoc < 0 || bpPowerSoc > 100 {
		return nil, errors.New("bpPowerSoc out of range. Valid range 0:100")
	}
	if minDsgSoc < 0 || minDsgSoc > 100 {
		return nil, errors.New("minDsgSoc out of range. Valid range 0:100")
	}
	if minChgSoc < 0 || minChgSoc > 100 {
		return nil, errors.New("minChgSoc out of range. Valid range 0:100")
	}
	params := make(map[string]interface{})
	params["isConfig"] = enabled
	params["bpPowerSoc"] = bpPowerSoc
	params["minDsgSoc"] = minDsgSoc
	params["minChgSoc"] = minChgSoc
	return s.setParameter(ctx, "watthConfig", PowerStationModuleTypePd, params)
}

// SetAcAlwaysOn Set AC always on (acAutoOutConfig: 0: disabled; 1: enabled;minAcOutSoc: minimum SoC for turning on "AC always on" )
// { "id":123456789, "version":"1.0", "sn":"R331ZEB4ZEAL0528", "moduleType":1, "operateType":"acAutoOutConfig", "params":{ "acAutoOutConfig":0, "minAcOutSoc":20 } }
func (s *PowerStation) SetAcAlwaysOn(ctx context.Context, enabled SettingSwitcher, minAcOutSoc int) (*CmdSetResponse, error) {
	if minAcOutSoc < 0 || minAcOutSoc > 100 {
		return nil, errors.New("minAcOutSoc out of range. Valid range 0:100")
	}
	params := make(map[string]interface{})
	params["acAutoOutConfig"] = enabled
	params["minAcOutSoc"] = minAcOutSoc
	return s.setParameter(ctx, "acAutoOutConfig", PowerStationModuleTypePd, params)
}

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
