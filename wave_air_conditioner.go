package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type WaveAirConditioner struct {
	c  *Client
	sn string
}

type ConditionerMainMode int

const (
	ConditionerMainModeCool ConditionerMainMode = 0
	ConditionerMainModeHeat ConditionerMainMode = 1
	ConditionerMainModeFan  ConditionerMainMode = 2
)

type ConditionerSubMode int

const (
	ConditionerSubModeMax    ConditionerSubMode = 0
	ConditionerSubModeSleep  ConditionerSubMode = 1
	ConditionerSubModeEco    ConditionerSubMode = 2
	ConditionerSubModeManual ConditionerSubMode = 3
)

type ConditionerTemperatureDisplayMode int

const (
	ConditionerTemperatureDisplayModeAmbient   ConditionerTemperatureDisplayMode = 0
	ConditionerTemperatureDisplayModeAirOutlet ConditionerTemperatureDisplayMode = 1
)

type ConditionerWindSpeed int

const (
	ConditionerWindSpeedLow    ConditionerWindSpeed = 0
	ConditionerWindSpeedMedium ConditionerWindSpeed = 1
	ConditionerWindSpeedHigh   ConditionerWindSpeed = 2
)

type ConditionerLightStripMode int

const (
	ConditionerLightStripModeFollowScreen ConditionerLightStripMode = 0
	ConditionerLightStripModeAlwaysOn     ConditionerLightStripMode = 1
	ConditionerLightStripModeAlwaysOff    ConditionerLightStripMode = 2
)

type ConditionerPowerMode int

const (
	ConditionerPowerModeStartup  ConditionerPowerMode = 1
	ConditionerPowerModeStandby  ConditionerPowerMode = 2
	ConditionerPowerModeShutdown ConditionerPowerMode = 3
)

func (c *WaveAirConditioner) GetSn() string {
	return c.sn
}

// SetMainMode Set main mode(0: Cool, 1: Heat, 2: Fan)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "moduleType":1, "operateType":"mainMode", "params":{ "mainMode":1 } }
func (c *WaveAirConditioner) SetMainMode(ctx context.Context, mainMode ConditionerMainMode) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["mainMode"] = mainMode
	return c.setParameter(ctx, "mainMode", ModuleTypePd, params)
}

// SetSubMode Set sub-mode(0: Max, 1: Sleep, 2: Eco, 3: Manual)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"subMode", "params":{ "subMode":3 } }
func (c *WaveAirConditioner) SetSubMode(ctx context.Context, subMode ConditionerSubMode) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["subMode"] = subMode
	return c.setParameter(ctx, "subMode", -1, params)
}

// SetTemperatureUnit Set unit of temperature(0: Celsius, 1: Fahrenheit)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"tempSys", "params":{ "mode":1 } }
func (c *WaveAirConditioner) SetTemperatureUnit(ctx context.Context, mode TemperatureUnit) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["mode"] = mode
	return c.setParameter(ctx, "tempSys", -1, params)
}

// SetScreenTimeout Set screen timeout (time unit: sec; Always on: "idleTime": 0, "idleMode": 0)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"display", "params":{ "idleTime":5, "idleMode":1 } }
func (c *WaveAirConditioner) SetScreenTimeout(ctx context.Context, idleTime int, hasScreenTimeout SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["idleTime"] = idleTime
	params["idleMode"] = hasScreenTimeout
	return c.setParameter(ctx, "display", -1, params)
}

// SetTimer Set timer(timeSet: 0-65535; Unit: min;timeEn: 0: Turn off 1: Turn on)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"sacTiming", "params":{ "timeSet":10, "timeEn":1 } }
func (c *WaveAirConditioner) SetTimer(ctx context.Context, timeSet int, timeEn SettingSwitcher) (*CmdSetResponse, error) {
	if timeSet < 0 || timeSet > 65535 {
		return nil, errors.New("timeSet is out of range. Range 0:65535")
	}
	params := make(map[string]interface{})
	params["timeSet"] = timeSet
	params["timeEn"] = timeEn
	return c.setParameter(ctx, "sacTiming", -1, params)
}

// SetEnableBuzzer Enable buzzer (0: Disable; 1: Enable)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"beepEn", "params":{ "en":1 } }
func (c *WaveAirConditioner) SetEnableBuzzer(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["en"] = enabled
	return c.setParameter(ctx, "beepEn", -1, params)
}

// SetTemperature Set temperature(16-30 ℃）
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "moduleType":1, "operateType":"setTemp", "params":{ "setTemp":27 } }
func (c *WaveAirConditioner) SetTemperature(ctx context.Context, setTemp int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["setTemp"] = setTemp
	return c.setParameter(ctx, "setTemp", ModuleTypePd, params)
}

// SetTemperatureDisplay Set temperature display (0: Display ambient temperature; 1: Display air outlet temperature)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "moduleType":1, "operateType":"tempDisplay", "params":{ "tempDisplay":0 } }
func (c *WaveAirConditioner) SetTemperatureDisplay(ctx context.Context, tempDisplay ConditionerTemperatureDisplayMode) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["tempDisplay"] = tempDisplay
	return c.setParameter(ctx, "setTemp", ModuleTypePd, params)
}

// SetWindSpeed Set wind speed (0: Low; 1: Medium; 2: High)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"fanValue", "params":{ "fanValue":1 } }
func (c *WaveAirConditioner) SetWindSpeed(ctx context.Context, fanValue ConditionerWindSpeed) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["fanValue"] = fanValue
	return c.setParameter(ctx, "fanValue", -1, params)
}

// SetAutomaticDrainage Set automatic drainage(
// In Cool/Fan mode: 0: Turn on Manual drainage，1: Turn on No drainage, 2: Turn off Manual drainage, 3 Turn off No drainage
// In Heat Mode: 0: Turn off, 1: Turn on Manual drainage， 3: Turn off Manual drainage)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "operateType":"wteFthEn", "params":{ "wteFthEn":3 } }
func (c *WaveAirConditioner) SetAutomaticDrainage(ctx context.Context, wteFthEn int) (*CmdSetResponse, error) {
	if wteFthEn < 0 || wteFthEn > 3 {
		return nil, errors.New("wteFthEn is out of range. Range 0:3")
	}
	params := make(map[string]interface{})
	params["wteFthEn"] = wteFthEn

	return c.setParameter(ctx, "wteFthEn", -1, params)
}

// SetLightStripMode Light strip settings (0: Follow the screen; 1: Always on; 2: Always off; other parameters indicate “Always off”)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "moduleType":1, "operateType":"rgbState", "params":{ "rgbState":1 } }
func (c *WaveAirConditioner) SetLightStripMode(ctx context.Context, rgbState ConditionerLightStripMode) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["rgbState"] = rgbState

	return c.setParameter(ctx, "rgbState", ModuleTypePd, params)
}

// SetPowerMode Remote startup/shutdown (1: Startup; 2: Standby; 3: Shutdown)
// { "id":123456789, "version":"1.0", "sn":"KT21ZCH2ZF170012", "moduleType":1, "operateType":"powerMode", "params":{ "powerMode":2 } }
func (c *WaveAirConditioner) SetPowerMode(ctx context.Context, powerMode ConditionerPowerMode) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["powerMode"] = powerMode

	return c.setParameter(ctx, "powerMode", ModuleTypePd, params)
}

func (c *WaveAirConditioner) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return c.c.GetDeviceParameters(ctx, c.sn, params)
}

func (c *WaveAirConditioner) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return c.c.GetDeviceAllParameters(ctx, c.sn)
}

func (c *WaveAirConditioner) setParameter(ctx context.Context, opType string, modType ModuleType, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:          fmt.Sprint(time.Now().UnixMilli()),
		OperateType: opType,
		Sn:          c.sn,
		Params:      params,
	}

	if modType != -1 {
		cmdReq.ModuleType = modType
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
	return c.c.SetDeviceParameter(ctx, req)
}
