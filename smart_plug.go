package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Ecoflow documentation: https://developer-eu.ecoflow.com/us/document/smartPlug

type SmartPlug struct {
	c  *Client
	sn string
}

func (s *SmartPlug) GetSn() string {
	return s.sn
}

// SetRelaySwitch Relay switch(0: off, 1: on)
// {"sn": "HW52ZDH1RF3J0033","cmdCode": "WN511_SOCKET_SET_PLUG_SWITCH_MESSAGE","params": {"plugSwitch": 0}}
func (s *SmartPlug) SetRelaySwitch(ctx context.Context, enabled SettingSwitcher) (*CmdSetResponse, error) {
	params := make(map[string]interface{})
	params["plugSwitch"] = enabled
	return s.setParameter(ctx, "WN511_SOCKET_SET_PLUG_SWITCH_MESSAGE", params)
}

// SetIndicatorBrightness Indicator light brightness adjustment(rgb brightness: 0-1023 (the larger the value, the higher the brightness); default value: 1023)
// {"sn": "HW52ZDH1RF3J0033","cmdCode": "WN511_SOCKET_SET_BRIGHTNESS_PACK","params": {"brightness": 1000}}
func (s *SmartPlug) SetIndicatorBrightness(ctx context.Context, brightness int) (*CmdSetResponse, error) {
	if brightness < 0 || brightness > 1023 {
		return nil, errors.New("brightness out of range. Expected value from 0 to 1023")
	}
	params := make(map[string]interface{})
	params["brightness"] = brightness
	return s.setParameter(ctx, "WN511_SOCKET_SET_BRIGHTNESS_PACK", params)
}

// DeleteScheduledTasks Deleting scheduled tasks(taskIndex: 0-9)
// {"sn": "HW52ZDH1RF3J0033","cmdCode": "WN511_SOCKET_DELETE_TIME_TASK","params": {"taskIndex": 1}}
func (s *SmartPlug) DeleteScheduledTasks(ctx context.Context, taskIndex int) (*CmdSetResponse, error) {
	if taskIndex < 0 || taskIndex > 9 {
		return nil, errors.New("taskIndex out of range. Expected value from 0 to 9")
	}
	params := make(map[string]interface{})
	params["taskIndex"] = taskIndex
	return s.setParameter(ctx, "WN511_SOCKET_DELETE_TIME_TASK", params)
}

func (s *SmartPlug) setParameter(ctx context.Context, cmdCode string, params map[string]interface{}) (*CmdSetResponse, error) {
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

func (s *SmartPlug) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

func (s *SmartPlug) GetAllParameters(ctx context.Context) (map[string]interface{}, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}
