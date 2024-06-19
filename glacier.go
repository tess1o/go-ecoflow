package ecoflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Ecoflow documentation:
// https://developer-eu.ecoflow.com/us/document/glacier

type Glacier struct {
	c  *Client
	sn string
}

func (g *Glacier) GetSn() string {
	return g.sn
}

type GlacierModeType int

const (
	GlacierModeTypeNormal GlacierModeType = 0
	GlacierModeTypeEco    GlacierModeType = 1
)

type GlacierIceShape int

const (
	GlacierIceShapeSmall GlacierIceShape = 0
	GlacierIceShapeLarge GlacierIceShape = 1
)

type GlacierSensorDetection int

const (
	GlacierSensorDetectionUnblocking GlacierSensorDetection = 0
	GlacierSensorDetectionBlocking   GlacierSensorDetection = 1
)

type GlacierBuzzerCommand int

const (
	GlacierBuzzerCommandAlwaysBeeping GlacierBuzzerCommand = 0
	GlacierBuzzerCommandBeepOnce      GlacierBuzzerCommand = 1
	GlacierBuzzerCommandBeepTwice     GlacierBuzzerCommand = 2
	GlacierBuzzerCommandThreeTimes    GlacierBuzzerCommand = 3
)

type GlacierVoltageProtectionLevel int

const (
	GlacierVoltageProtectionLevelLow    GlacierVoltageProtectionLevel = 0
	GlacierVoltageProtectionLevelMedium GlacierVoltageProtectionLevel = 1
	GlacierVoltageProtectionLevelHigh   GlacierVoltageProtectionLevel = 2
)

// SetTemperature Set temperature(tmpR indicates the temperature of the right side of the refrigerator,
// tmpL indicates the temperature of the left side, and tmpM indicates the temperature setting after the middle partition is removed.
// The difference between tmpR and tmpL cannot exceed 25℃)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"temp", "params":{ "tmpR":-19, "tmpL":0, "tmpM":0 } }
func (g *Glacier) SetTemperature(ctx context.Context, tmpR, tmpL, tmpM int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["tmpR"] = tmpR
	params["tmpL"] = tmpL
	params["tmpM"] = tmpM
	return g.setParameter(ctx, "temp", params)
}

// SetEcoMode Set ECO mode(mode: 1: ECO; 0: Normal)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"ecoMode", "params":{ "mode":1 } }
func (g *Glacier) SetEcoMode(ctx context.Context, mode GlacierModeType) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["mode"] = mode
	return g.setParameter(ctx, "ecoMode", params)
}

// SetBuzzerEnablingStatus Set buzzer enabling status(0: Disable; 1: Enable)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"beepEn", "params":{ "flag":1 } }
func (g *Glacier) SetBuzzerEnablingStatus(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["flag"] = enabled
	return g.setParameter(ctx, "beepEn", params)
}

// SetBuzzerCommand Buzzer commands(1: Beep once; 2: Beep twice; 3: Beep three times; 0: Always beeping)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"beep", "params":{ "flag":1 } }
func (g *Glacier) SetBuzzerCommand(ctx context.Context, command GlacierBuzzerCommand) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["flag"] = command
	return g.setParameter(ctx, "beep", params)
}

// SetScreenTimeout Set screen timeout(unit: sec; when set to 0, the screen is always on)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"blTime", "params":{ "time":600 } }
func (g *Glacier) SetScreenTimeout(ctx context.Context, time int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["time"] = time
	return g.setParameter(ctx, "blTime", params)
}

// SetTemperatureUnit Set temperature unit(0: Celsius; 1: Fahrenheit）
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"tmpUnit", "params":{ "unit":0 } }
func (g *Glacier) SetTemperatureUnit(ctx context.Context, unit TemperatureUnit) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["unit"] = unit
	return g.setParameter(ctx, "tmpUnit", params)
}

// SetIceMaking Set ice making(If "enable"=0, ice making is disabled.
// If "enable"=1 and "iceShape"=0, the device will make small ice cubes. If "enable"=1 and "iceShape"=1, the device will make large ice cubes.)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"iceMake", "params":{ "enable":1, "iceShape":1 } }
func (g *Glacier) SetIceMaking(ctx context.Context, enable SettingSwitcher, iceShape GlacierIceShape) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["enable"] = enable
	params["iceShape"] = iceShape
	return g.setParameter(ctx, "iceMake", params)
}

// SetIceDetaching Set ice detaching(enable: 0: Invalid, 1: Detach iceiceTm: Duration of ice detaching; unit: secfsmState: 4: Detaching ice, 5: Detaching completed）
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"deIce", "params":{ "enable":0 } }
func (g *Glacier) SetIceDetaching(ctx context.Context, enable SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["enable"] = enable
	return g.setParameter(ctx, "deIce", params)
}

// SetSensorDetectionBlocking Sensor detection blocking(0: Unblocked; 1: Blocked)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"sensorAdv", "params":{ "sensorAdv":1 } }
func (g *Glacier) SetSensorDetectionBlocking(ctx context.Context, sensor GlacierSensorDetection) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["sensorAdv"] = sensor
	return g.setParameter(ctx, "sensorAdv", params)
}

// SetBatteryLowVoltageProtectionLevel Set battery low voltage protection level(state: 0: Disabled; 1: Enabled; level: 0: Low; 1: Medium; 2: High)
// { "id":123456789, "version":"1.0", "sn":"BX11ZCB4EF2E0002", "moduleType":1, "operateType":"protectBat", "params":{ "state":1, "level":0 } }
func (g *Glacier) SetBatteryLowVoltageProtectionLevel(ctx context.Context, state SettingSwitcher, level GlacierVoltageProtectionLevel) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["state"] = state
	params["level"] = level
	return g.setParameter(ctx, "protectBat", params)
}

// GetParameter get specified parameters for the device
func (g *Glacier) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return g.c.GetDeviceParameters(ctx, g.sn, params)
}

// GetAllParameters get all parameters for the device
func (g *Glacier) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return g.c.GetDeviceAllParameters(ctx, g.sn)
}

// internal function to generate a request for setting the parameters
func (g *Glacier) setParameter(ctx context.Context, opType string, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:          fmt.Sprint(time.Now().UnixMilli()),
		OperateType: opType,
		ModuleType:  ModuleTypePd,
		Sn:          g.sn,
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
	return g.c.SetDeviceParameter(ctx, req)
}
