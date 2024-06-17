package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Ecoflow documentation: https://developer-eu.ecoflow.com/us/document/shp

type SmartHomePanel struct {
	c  *Client
	sn string
}

func (s *SmartHomePanel) GetSn() string {
	return s.sn
}

// SetRtcTime RTC time update
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 3, "week": 2, "sec": 17, "min": 38, "hour": 18, "day": 16, "month": 11, "year": 2022 } }
func (s *SmartHomePanel) SetRtcTime(ctx context.Context, t time.Time) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 3
	params["week"] = int(t.Weekday())
	params["sec"] = t.Second()
	params["min"] = t.Minute()
	params["hour"] = t.Hour()
	params["day"] = t.Day()
	params["month"] = int(t.Month())
	params["year"] = t.Year()

	return s.setParameter(ctx, params)
}

// SetLoadChannelControl Load channel control
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 16, "ch": 1, "ctrlMode": 1, "sta": 1 } }
func (s *SmartHomePanel) SetLoadChannelControl(ctx context.Context, ch, ctrlMode, sta int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 16
	params["ch"] = ch
	params["ctrlMode"] = ctrlMode
	params["sta"] = sta

	return s.setParameter(ctx, params)
}

// SetStandByChannelControl Standby channel control
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 17, "ch": 10, "ctrlMode": 1, "sta": 1 } }
func (s *SmartHomePanel) SetStandByChannelControl(ctx context.Context, ch, ctrlMode, sta int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 17
	params["ch"] = ch
	params["ctrlMode"] = ctrlMode
	params["sta"] = sta

	return s.setParameter(ctx, params)
}

//TODO: Split-phase information configuration
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 18, "cfgList": [ { "linkMark": 1, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 }, { "linkMark": 0, "linkCh": 0 } ] } }

// SetChannelCurrentConfiguration Channel current configuration (cur: 6, 13, 16, 20, 30)
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 20, "chNum": 0, "cur": 6 } }
func (s *SmartHomePanel) SetChannelCurrentConfiguration(ctx context.Context, chNum, cur int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 20
	params["chNum"] = chNum
	params["cur"] = cur

	return s.setParameter(ctx, params)
}

// SetGridPowerConfiguration Grid power parameter configuration (gridVol: 220 230 240)
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "gridVol": 230, "gridFreq": 50, "cmdSet": 11, "id": 22 } }
func (s *SmartHomePanel) SetGridPowerConfiguration(ctx context.Context, gridVol, gridFreq int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 22
	params["gridVol"] = gridVol
	params["gridFreq"] = gridFreq

	return s.setParameter(ctx, params)
}

// SetEspMode EPS mode configuration (eps: 0: off, 1: on)
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 24, "eps": 1 } }
func (s *SmartHomePanel) SetEspMode(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 24
	params["eps"] = enabled

	return s.setParameter(ctx, params)
}

// SetChannelEnableStatusConfiguration Channel enable status configuration(chNum: 0–9, isEnable, 0: off, 1: on)
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "isEnable": 1, "chNum": 1, "cmdSet": 11, "id": 26 } }
func (s *SmartHomePanel) SetChannelEnableStatusConfiguration(ctx context.Context, chNum int, enabled SettingSwitcher) (*CmdSetResponse, error) {
	if chNum < 0 || chNum > 9 {
		return nil, errors.New("chNum is out of range. Range 0:9")
	}

	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 24
	params["chNum"] = chNum
	params["isEnable"] = enabled

	return s.setParameter(ctx, params)
}

// SetLoadChannelConfiguration Load channel information configuration(chNum 0~9 )
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 32, "chNum": 1, "info": { "chName": "test", "iconInfo": 10 } } }
func (s *SmartHomePanel) SetLoadChannelConfiguration(ctx context.Context, chNum int, chName string, iconInfo int) (*CmdSetResponse, error) {
	if chNum < 0 || chNum > 9 {
		return nil, errors.New("chNum is out of range. Range 0:9")
	}

	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 32
	params["chNum"] = chNum

	info := make(map[string]interface{})
	info["chName"] = chName
	info["iconInfo"] = iconInfo

	params["info"] = info

	return s.setParameter(ctx, params)
}

// SetRegionInformation Region information configuration
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 34, "area": "US, China" } }
func (s *SmartHomePanel) SetRegionInformation(ctx context.Context, area string) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 34
	params["area"] = area

	return s.setParameter(ctx, params)
}

//TODO: Setting the emergency mode

//TODO: Setting the scheduled charging job

//TODO: Setting the scheduled discharging job

// SetConfigurationStatus Setting the configuration status
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 7, "cfgSta": 1 } }
func (s *SmartHomePanel) SetConfigurationStatus(ctx context.Context, cfgSta SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 7
	params["cfgSta"] = cfgSta

	return s.setParameter(ctx, params)
}

// StartSelfCheckInformationPushing Start self-check information pushing
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 112, "selfCheckType": 1 } }
func (s *SmartHomePanel) StartSelfCheckInformationPushing(ctx context.Context, selfCheckType int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 112
	params["selfCheckType"] = selfCheckType

	return s.setParameter(ctx, params)
}

// PushStandByChargingDischargingParameters Pushing standby charging/discharging parameters
// { "sn": "SP10ZAW5ZE9E0052", "operateType": "TCP", "params": { "cmdSet": 11, "id": 29, "forceChargeHigh": 0, "discLower": 0 } }
func (s *SmartHomePanel) PushStandByChargingDischargingParameters(ctx context.Context, forceChargeHigh, discLower int) (*CmdSetResponse, error) {
	params := make(map[string]interface{})

	params["cmdSet"] = 11
	params["id"] = 29
	params["forceChargeHigh"] = forceChargeHigh
	params["discLower"] = discLower

	return s.setParameter(ctx, params)
}

func (s *SmartHomePanel) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

func (s *SmartHomePanel) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

func (s *SmartHomePanel) setParameter(ctx context.Context, params map[string]interface{}) (*CmdSetResponse, error) {
	cmdReq := CmdSetRequest{
		Id:          fmt.Sprint(time.Now().UnixMilli()),
		OperateType: "TCP",
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
