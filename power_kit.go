package ecoflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PowerKit struct {
	c        *Client
	sn       string
	moduleSn string
}

func (k *PowerKit) GetSn() string {
	return k.sn
}

type powerKitSetCmdRequest struct {
	Id          string                 `json:"id"`
	Sn          string                 `json:"sn"`
	ModuleSn    string                 `json:"moduleSn"`
	ModuleType  PowerKitModuleType     `json:"moduleType"`
	OperateType string                 `json:"operateType"`
	Params      map[string]interface{} `json:"params"`
}

type PowerKitModuleType int

//15362: BBC_IN
//15363: BBC_OUT
//15365: IC_LOW
//0: BP5000/BP2000
//15367: LD_AC
//15368: LD_DC
//15370: Wireless

const (
	PowerKitModuleTypeBbcIn        PowerKitModuleType = 15362
	PowerKitModuleTypeBbcOut       PowerKitModuleType = 15363
	PowerKitModuleTypeIcLow        PowerKitModuleType = 15365
	PowerKitModuleTypeBp5000Bp2000 PowerKitModuleType = 0
	PowerKitModuleTypeLdAc         PowerKitModuleType = 15367
	PowerKitModuleTypeLdDc         PowerKitModuleType = 15368
	PowerKitModuleTypeWireless     PowerKitModuleType = 15370
)

type PowerKitDcVoltage int

const (
	PowerKitDcVoltage12V PowerKitDcVoltage = 0
	PowerKitDcVoltage24V PowerKitDcVoltage = 1
)

// BBC_IN

// SetDcOutputVoltage DC output voltage (0: 12 V, 1: 24 V)
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M1093-DCIN-CA7C3", "moduleType": 15362, "operateType": "dischgParaSet", "params": { "volTag": 0 } }
func (k *PowerKit) SetDcOutputVoltage(ctx context.Context, voltage PowerKitDcVoltage) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["volTag"] = voltage
	return k.setParameter(ctx, "quietMode", PowerKitModuleTypeBbcIn, params)
}

// SetChargingSettings Charging settings
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M1093-DCIN-CA7C3", "moduleType": 15362,
// "operateType": "chgParaSet", "params": { "chgPause": 0, "maxChgCurr": 30, "altVoltLmtEn": 255, "shakeCtrlDisable": 255,
// "altCableUnit": 255, "altCableLen": -1, "altVoltLmt": 65535 } }
func (k *PowerKit) SetChargingSettings(ctx context.Context, chgPause, maxChgCurr, altVoltLmtEn, shakeCtrlDisable, altCableUnit, altCableLen, altVoltLmt int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["chgPause"] = chgPause
	params["maxChgCurr"] = maxChgCurr
	params["altVoltLmtEn"] = altVoltLmtEn
	params["shakeCtrlDisable"] = shakeCtrlDisable
	params["altCableUnit"] = altCableUnit
	params["altCableLen"] = altCableLen
	params["altVoltLmt"] = altVoltLmt
	return k.setParameter(ctx, "chgParaSet", PowerKitModuleTypeBbcIn, params)
}

//BBC_OUT

// SetDischargingSettings Discharging settings(swSta: 0: off 1: on)
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M1093-DCIN-CA7C3", "moduleType": 15362, "operateType": "dischgParaSet", "params": { "swSta": 0 } }
// in documentation it says BBC_OUT (15363), in the json example it's still 15362. Next examples uses 15363...
func (k *PowerKit) SetDischargingSettings(ctx context.Context, enabled PowerKitDcVoltage) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["swSta"] = enabled
	return k.setParameter(ctx, "dischgParaSet", PowerKitModuleTypeBbcOut, params)
}

// SetBroadcastInstructionForRTCTimeSynchronization A broadcast instruction for synchronizing RTC time
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M109ZEB4Z0000016", "moduleType": 15363, "operateType": "rtcBroadcast",
// "params": { "unixTime": 1710835118, "timeZone": 8, "timeZoneQuarter": 1 } }
func (k *PowerKit) SetBroadcastInstructionForRTCTimeSynchronization(ctx context.Context, unixTime int64, timeZone int, timeZoneQuarter int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["unixTime"] = unixTime
	params["timeZone"] = timeZone
	params["timeZoneQuarter"] = timeZoneQuarter
	return k.setParameter(ctx, "rtcBroadcast", PowerKitModuleTypeBbcOut, params)
}

// IC_LOW

// SetCommandForDischarging Command for discharging, powerOn: 0: AC off, 1: AC on
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M109ZEB4ZE7B0963", "moduleType": 15365, "operateType": "dischgIcParaSet",
// "params": { "acCurrMaxSet": 255, "powerOn": 0, "acChgDisa": 255, "acFrequencySet": 255, "acVolSet": 255 } }
func (k *PowerKit) SetCommandForDischarging(ctx context.Context, acCurrMaxSet int, powerOn SettingSwitcher, acChgDisa, acFrequencySet, acVolSet int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["acCurrMaxSet"] = acCurrMaxSet
	params["powerOn"] = powerOn
	params["acChgDisa"] = acChgDisa
	params["acFrequencySet"] = acFrequencySet
	params["acVolSet"] = acVolSet
	return k.setParameter(ctx, "dischgIcParaSet", PowerKitModuleTypeIcLow, params)
}

// SetAcInputCurrent AC input current (range: 1-23)
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M109ZEB4ZE7B0963", "moduleType": 15365, "operateType": "dischgIcParaSet", "params": { "acCurrMaxSet": 10 } }
func (k *PowerKit) SetAcInputCurrent(ctx context.Context, acCurrMaxSet int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["acCurrMaxSet"] = acCurrMaxSet
	return k.setParameter(ctx, "dischgIcParaSet", PowerKitModuleTypeIcLow, params)
}

// SetGridPowerInPriority Grid power in priority (passByModeEn, 1: on, 2: off)
// { "id": 123456789, "version": "1.0", "sn": "M106ZAB4Z000001F", "moduleSn": "M109ZEB4ZE7B0963", "moduleType": 15365, "operateType": "dsgIcParaSet",
// "params": { "dsgLowPwrEn": 255, "pfcDsgModeEn": 255, "passByCurrMax": 255, "passByModeEn": 1 } }
func (k *PowerKit) SetGridPowerInPriority(ctx context.Context, dsgLowPwrEn, pfcDsgModeEn, passByCurrMax int, passByModeEn SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["dsgLowPwrEn"] = dsgLowPwrEn
	params["pfcDsgModeEn"] = pfcDsgModeEn
	params["passByCurrMax"] = passByCurrMax
	params["passByModeEn"] = passByModeEn
	return k.setParameter(ctx, "dsgIcParaSet", PowerKitModuleTypeIcLow, params)
}

func (k *PowerKit) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return k.c.GetDeviceParameters(ctx, k.sn, params)
}

func (k *PowerKit) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return k.c.GetDeviceAllParameters(ctx, k.sn)
}

func (k *PowerKit) setParameter(ctx context.Context, opType string, modType PowerKitModuleType, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := powerKitSetCmdRequest{
		Id:          fmt.Sprint(time.Now().UnixMilli()),
		Sn:          k.sn,
		ModuleSn:    k.moduleSn,
		ModuleType:  modType,
		OperateType: opType,
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
	return k.c.SetDeviceParameter(ctx, req)
}
