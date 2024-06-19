# Ecoflow API implementation via REST API in Go

Ecoflow Rest API implementation in Go (no external dependencies) that allows to get list of linked devices, get all
devices parameters and set devices settings. \
The library was tested on Ecoflow Delta 2 and Ecoflow River 2 (I don't have other their products)\
The Ecoflow documentation is not complete and sometimes hard to understand, thus some APIs might not work as expected

## Installation
To get the library just use `go get` command:
`go get github.com/tess1o/go-ecoflow`

## Supported devices:
1. Power Stations (regular ecoflow power stations, like Delta 2, River 2, etc)
2. Power Stations (PRO)
3. Smart Plug
4. PowerStream Micro Inverter
5. Smart Home Panel
6. Wave Air Conditioner
7. Glacier

## Features

The library allows to:

1. Get list of all linked devices
2. Get specified parameters from a device
3. Get all parameters from a device
4. Change device's settings

Basically that's all documented features the Ecoflow REST API provides

## Documentation

Link to official documentation: https://developer-eu.ecoflow.com/us/document/introduction

## How to get Access Token and Secret Token

1. Go to https://developer-eu.ecoflow.com/
2. Click on "Become a Developer"
3. Login with your Ecoflow username and Password
4. Wait until the access is approved by Ecoflow
5. Receive email with subject "Approval notice from EcoFlow Developer Platform". May take some time
6. Go to https://developer-eu.ecoflow.com/us/security and create new AccessKey and SecretKey

## Usage example

Usage example. See more details for each device below in the README.

```go
package main

import (
	"context"
	"github.com/tess1o/go-ecoflow"
	"log/slog"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		slog.Error("AccessKey and SecretKey are mandatory")
		return
	}

	//create new client.
	client := ecoflow.NewEcoflowClient(accessKey, secretKey)

	// creating new client with options. Current supports two options:
	// 1. custom ecoflow base url (can be used with proxies, or if they change the url)
	// 2. custom http client

	// client = ecoflow.NewEcoflowClient(accessKey, secretKey,
	//	ecoflow.WithBaseUrl("https://ecoflow-api.example.com"),
	//	ecoflow.WithHttpClient(customHttpClient()),
	//)

	//get all linked ecoflow devices. Returns SN and online status
	client.GetDeviceList(context.Background())
	// get param1 and param2 for device 
	client.GetDeviceParameters(context.Background(), "DEVICE_SERIAL_NUMBER", []string{"param1", "param2"})
	// get all parameters for device
	client.GetDeviceAllParameters(context.Background(), "DEVICE_SERIAL_NUMBER")

	ctx := context.Background()

	// get set / get functions for power stations. PRO version is not currently implemented
	ps := client.GetPowerStation("DEVICE_SERIAL_NUMBER")

	//set functions, see details in documentation to each function. There are much more functions for each type of device
	ps.SetDcSwitch(ctx, ecoflow.SettingEnabled)
	ps.SetMaxChargeSoC(ctx, 99)
	ps.SetMinDischargeSoC(ctx, 1)
	ps.SetStandByTime(ctx, 60)
	ps.SetCarStandByTime(ctx, 60)
	ps.SetPrioritizePolarCharging(ctx, ecoflow.SettingEnabled)

	//get parameters functions
	ps.GetAllParameters(ctx)
	ps.GetParameter(ctx, []string{"mppt.acStandbyMins", "mppt.dcChgCurrent"})
}

```

## Library API description

Each function has documentation taken from the Ecoflow website. I don't have more information about the API.

### Client

`ecoflow.Client` must be created using the `accessToken` and `secretToken`.
It provides function to get all linked devices and generic functions to get all parameters for any Ecoflow device or set parameters for any
Device. It can be useful if a new devices is introduced and not supported by this library or new APIs are introduced and
not implemented by the library. Usually you don't need to use `client.GetDeviceAllParameters` or `client.SetParameters`,
it's implemented in each device type.\
Additionally `ecoflow.Client` provides functions to get API for each type of Ecoflow device (see examples below)

```go
//get client by providing access key and secret key
client := ecoflow.NewEcoflowClient(accessKey, secretKey)

// creating new client with options. Current supports two options:
// 1. custom ecoflow base url (can be used with proxies, or if they change the url)
// 2. custom http client

// client = ecoflow.NewEcoflowClient(accessKey, secretKey,
//	ecoflow.WithBaseUrl("https://ecoflow-api.example.com"),
//	ecoflow.WithHttpClient(customHttpClient()),
//)

//get all linked ecoflow devices. Returns SN and online status
client.GetDeviceList(context.Background())
// get param1 and param2 for device 
client.GetDeviceParameters(context.Background(), "DEVICE_SERIAL_NUMBER", []string{"param1", "param2"})
// get all parameters for device
client.GetDeviceAllParameters(context.Background(), "DEVICE_SERIAL_NUMBER")

// set parameters, where params is the map[string]interface{}.
client.SetParameters(context.Background(), params)
```

### Power Station

API that can be used with an Ecoflow PowerStation (not PRO) version.

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetPowerStation("POWER_STATION_SERIAL_NUMBER")
```

The list of available functions:

```
func (s *PowerStation) GetSn()(string)

func (s *PowerStation) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)
func (s *PowerStation) GetAllParameters(ctx context.Context)

func (s *PowerStation) SetBuzzerSilentMode(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStation) SetCarChargerSwitch(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStation) SetAcEnabled(ctx context.Context, acEnabled , xBoostEnabled SettingSwitcher, outFreq GridFrequency, outVoltage int)(*CmdSetResponse, error)
func (s *PowerStation) SetAcChargingSettings(ctx context.Context, chargeWatts int, chgPauseFlag SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStation) SetAcStandByTime(ctx context.Context, standbyMins int)(*CmdSetResponse, error)
func (s *PowerStation) SetCarStandByTime(ctx context.Context, standbyMins int)(*CmdSetResponse, error)
func (s *PowerStation) Set12VDcChargingCurrent(ctx context.Context, chargingCurrent int)(*CmdSetResponse, error)
func (s *PowerStation) SetPvChargingTypeSettings(ctx context.Context, chargeType1 , chargeType2 PowerStationPvChargeType)(*CmdSetResponse, error)
func (s *PowerStation) SetStandByTime(ctx context.Context, standbyMin int)(*CmdSetResponse, error)
func (s *PowerStation) SetDcSwitch(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStation) SetLcdScreenTimeout(ctx context.Context, delayOffSeconds int)(*CmdSetResponse, error)
func (s *PowerStation) SetPrioritizePolarCharging(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStation) SetEnergyManagement(ctx context.Context, enabled SettingSwitcher, bpPowerSoc , minDsgSoc , minChgSoc int)(*CmdSetResponse, error)
func (s *PowerStation) SetAcAlwaysOn(ctx context.Context, enabled SettingSwitcher, minAcOutSoc int)(*CmdSetResponse, error)
func (s *PowerStation) SetMaxChargeSoC(ctx context.Context, maxChgSoc int)(*CmdSetResponse, error)
func (s *PowerStation) SetMinDischargeSoC(ctx context.Context, minDsgSoc int)(*CmdSetResponse, error)
func (s *PowerStation) SetSoCToTurnOnSmartGenerator(ctx context.Context, openOilSoc int)(*CmdSetResponse, error)
func (s *PowerStation) SetSoCToTurnOffSmartGenerator(ctx context.Context, closeOilSoc int)(*CmdSetResponse, error)
```

### Power Station (PRO)

API that can be used with an Ecoflow Power Station PRO versions

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetPowerStationPro("POWER_STATION_PRO_SERIAL_NUMBER")
```

The list of available functions:

```
func (s *PowerStationPro) GetSn()(string)

func (s *PowerStationPro) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)
func (s *PowerStationPro) GetAllParameters(ctx context.Context)

func (s *PowerStationPro) SetXboostSwitcher(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStationPro) SetCarChargerSwitch(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStationPro) SetMaxChargeLevel(ctx context.Context, maxChgSoc int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetMinDischargeLevel(ctx context.Context, minDsgSoc int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetCarInputCurrent(ctx context.Context, currMa int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetBeepSwitch(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *PowerStationPro) SetScreenBrightness(ctx context.Context, lcdBrightness int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetSoCToTurnOnSmartGenerator(ctx context.Context, openOilSoc int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetSoCToTurnOffSmartGenerator(ctx context.Context, closeOilSoc int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetUnitTimeout(ctx context.Context, standByMode int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetScreenTimeout(ctx context.Context, lcdTime int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetAcStandByTime(ctx context.Context, standByMins int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetAcChargingSettings(ctx context.Context, slowChgPower int)(*CmdSetResponse, error)
func (s *PowerStationPro) SetPvChargingType(ctx context.Context, chgType PowerStationPvChargeType)(*CmdSetResponse, error)
func (s *PowerStationPro) SetBypassAcAutoStart(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
```

### Power Kits

API that can be used with an Ecoflow Power Kits

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetPowerKit("POWER_KIT_SERIAL_NUMBER")
```

The list of available functions:

```
func (k *PowerKit) GetSn()(string)

func (k *PowerKit) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)
func (k *PowerKit) GetAllParameters(ctx context.Context)

func (k *PowerKit) SetDcOutputVoltage(ctx context.Context, voltage PowerKitDcVoltage)(*CmdSetResponse, error)
func (k *PowerKit) SetChargingSettings(ctx context.Context, chgPause , maxChgCurr , altVoltLmtEn , shakeCtrlDisable , altCableUnit , altCableLen , altVoltLmt int)(*CmdSetResponse, error)
func (k *PowerKit) SetDischargingSettings(ctx context.Context, enabled PowerKitDcVoltage)(*CmdSetResponse, error)
func (k *PowerKit) SetBroadcastInstructionForRTCTimeSynchronization(ctx context.Context, unixTime int64, timeZone int, timeZoneQuarter int)(*CmdSetResponse, error)
func (k *PowerKit) SetCommandForDischarging(ctx context.Context, acCurrMaxSet int, powerOn SettingSwitcher, acChgDisa , acFrequencySet , acVolSet int)(*CmdSetResponse, error)
func (k *PowerKit) SetAcInputCurrent(ctx context.Context, acCurrMaxSet int)(*CmdSetResponse, error)
func (k *PowerKit) SetGridPowerInPriority(ctx context.Context, dsgLowPwrEn , pfcDsgModeEn , passByCurrMax , passByModeEn int)(*CmdSetResponse, error)
func (k *PowerKit) SetChargingUpperLimit(ctx context.Context, maxChgSoc int)(*CmdSetResponse, error)
func (k *PowerKit) SetDischargingLowerLimit(ctx context.Context, minDsgSoc int)(*CmdSetResponse, error)
func (k *PowerKit) SetScreenStandByTime(ctx context.Context, standByTimeMinutes int)(*CmdSetResponse, error)
func (k *PowerKit) SetBpOff(ctx context.Context, enable SettingSwitcher)(*CmdSetResponse, error)
func (k *PowerKit) SetHeatingByDischarging(ctx context.Context, enable SettingSwitcher)(*CmdSetResponse, error)
func (k *PowerKit) SetClearingChargingErrors(ctx context.Context, clear SettingSwitcher)(*CmdSetResponse, error)
func (k *PowerKit) SetLowerLimitForStartupGenerator(ctx context.Context, soc int)(*CmdSetResponse, error)
func (k *PowerKit) SetUpperLimitForStartupGenerator(ctx context.Context, soc int)(*CmdSetResponse, error)
func (k *PowerKit) SetSixWayChannelRelayStatus(ctx context.Context, bitsSwSta SettingSwitcher)(*CmdSetResponse, error)
func (k *PowerKit) SetProductName(ctx context.Context, nameLen int, name string)(*CmdSetResponse, error)
func (k *PowerKit) SetScenarios(ctx context.Context, scenes int)(*CmdSetResponse, error)
func (k *PowerKit) SetTriggeringComprehensiveDataReport(ctx context.Context, times int)(*CmdSetResponse, error)
func (k *PowerKit) SetOilPocketStart(ctx context.Context, bitsSwSta SettingSwitcher)(*CmdSetResponse, error)
```

### Power Stream Micro Inverter

API that can be used with an Ecoflow Power Stream Micro Inverter

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetPowerStreamMicroInverter("INVERTER_SERIAL_NUMBER")
```

The list of available functions:

```
func (s *PowerStreamMicroInverter) GetSn()(string)

func (s *PowerStreamMicroInverter) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)
func (s *PowerStreamMicroInverter) GetAllParameters(ctx context.Context)

func (s *PowerStreamMicroInverter) SetPowerSupplyPriority(ctx context.Context, supplyPriority int)(*CmdSetResponse, error)
func (s *PowerStreamMicroInverter) SetCustomLoadPowerSettings(ctx context.Context, permanentWatts float64)(*CmdSetResponse, error)
func (s *PowerStreamMicroInverter) SetLowerLimitSettingsForBatterDischarging(ctx context.Context, lowerLimit float64)(*CmdSetResponse, error)
func (s *PowerStreamMicroInverter) SetUpperLimitSettingsForBatterCharging(ctx context.Context, upperLimit float64)(*CmdSetResponse, error)
func (s *PowerStreamMicroInverter) SetLightBrightness(ctx context.Context, brightness float64)(*CmdSetResponse, error)
func (s *PowerStreamMicroInverter) DeleteScheduledSwitchingTasks(ctx context.Context, taskIndex float64)(*CmdSetResponse, error)
```

### Wave Air Conditioner

API that can be used with an Ecoflow Wave Air Conditioner

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetWaveAirConditioner("WAVE_CONDITIONER_SERIAL_NUMBER")
```

The list of available functions:

```
func (c *WaveAirConditioner) GetSn()(string)

func (c *WaveAirConditioner) GetAllParameters(ctx context.Context)
func (c *WaveAirConditioner) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)

func (c *WaveAirConditioner) SetMainMode(ctx context.Context, mainMode ConditionerMainMode)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetSubMode(ctx context.Context, subMode ConditionerSubMode)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetTemperatureUnit(ctx context.Context, mode TemperatureUnit)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetScreenTimeout(ctx context.Context, idleTime int, hasScreenTimeout SettingSwitcher)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetTimer(ctx context.Context, timeSet int, timeEn SettingSwitcher)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetEnableBuzzer(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetTemperature(ctx context.Context, setTemp int)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetTemperatureDisplay(ctx context.Context, tempDisplay ConditionerTemperatureDisplayMode)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetWindSpeed(ctx context.Context, fanValue ConditionerWindSpeed)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetAutomaticDrainage(ctx context.Context, wteFthEn int)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetLightStripMode(ctx context.Context, rgbState ConditionerLightStripMode)(*CmdSetResponse, error)
func (c *WaveAirConditioner) SetPowerMode(ctx context.Context, powerMode ConditionerPowerMode)(*CmdSetResponse, error)
```

### Smart Plug

API that can be used with an Ecoflow Smart Plug

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetSmartPlug("WAVE_CONDITIONER_SERIAL_NUMBER")
```

The list of available functions:

```
func (s *SmartPlug) GetSn()(string)

func (s *SmartPlug) GetAllParameters(ctx context.Context)
func (s *SmartPlug) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)

func (s *SmartPlug) SetRelaySwitch(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *SmartPlug) SetIndicatorBrightness(ctx context.Context, brightness int)(*CmdSetResponse, error)
func (s *SmartPlug) DeleteScheduledTasks(ctx context.Context, taskIndex int)(*CmdSetResponse, error)
```

### Smart Home Panel

API that can be used with a Smart Home Panel

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetSmartHomePanel("HOME_PANEL_SERIAL_NUMBER")
```

The list of available functions:

```
func (s *SmartHomePanel) GetSn()(string)

func (s *SmartHomePanel) GetAllParameters(ctx context.Context)
func (s *SmartHomePanel) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)

func (s *SmartHomePanel) SetRtcTime(ctx context.Context, t time.Time)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetLoadChannelControl(ctx context.Context, ch , ctrlMode , sta int)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetStandByChannelControl(ctx context.Context, ch , ctrlMode , sta int)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetChannelCurrentConfiguration(ctx context.Context, chNum , cur int)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetGridPowerConfiguration(ctx context.Context, gridVol , gridFreq int)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetEspMode(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetChannelEnableStatusConfiguration(ctx context.Context, chNum int, enabled SettingSwitcher)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetLoadChannelConfiguration(ctx context.Context, chNum int, chName string, iconInfo int)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetRegionInformation(ctx context.Context, area string)(*CmdSetResponse, error)
func (s *SmartHomePanel) SetConfigurationStatus(ctx context.Context, cfgSta SettingSwitcher)(*CmdSetResponse, error)
func (s *SmartHomePanel) StartSelfCheckInformationPushing(ctx context.Context, selfCheckType int)(*CmdSetResponse, error)
func (s *SmartHomePanel) PushStandByChargingDischargingParameters(ctx context.Context, forceChargeHigh , discLower int)(*CmdSetResponse, error)
```

### Glacier

API that can be used with a Glacier

```go
client := ecoflow.NewEcoflowClient(accessKey, secretKey)
device := client.GetClacier("GLACIER_SERIAL_NUMBER")
```

The list of available functions:

```
func (g *Glacier) GetSn()(string)

func (g *Glacier) GetAllParameters(ctx context.Context)func (g *Glacier) SetTemperature(ctx context.Context, tmpR , tmpL , tmpM int)(*CmdSetResponse, error)
func (g *Glacier) GetParameter(ctx context.Context, params []string)(*GetCmdResponse, error)

func (g *Glacier) SetEcoMode(ctx context.Context, mode GlacierModeType)(*CmdSetResponse, error)
func (g *Glacier) SetBuzzerEnablingStatus(ctx context.Context, enabled SettingSwitcher)(*CmdSetResponse, error)
func (g *Glacier) SetBuzzerCommand(ctx context.Context, command GlacierBuzzerCommand)(*CmdSetResponse, error)
func (g *Glacier) SetScreenTimeout(ctx context.Context, time int)(*CmdSetResponse, error)
func (g *Glacier) SetTemperatureUnit(ctx context.Context, unit TemperatureUnit)(*CmdSetResponse, error)
func (g *Glacier) SetIceMaking(ctx context.Context, enable SettingSwitcher, iceShape GlacierIceShape)(*CmdSetResponse, error)
func (g *Glacier) SetIceDetaching(ctx context.Context, enable SettingSwitcher)(*CmdSetResponse, error)
func (g *Glacier) SetSensorDetectionBlocking(ctx context.Context, sensor GlacierSensorDetection)(*CmdSetResponse, error)
func (g *Glacier) SetBatteryLowVoltageProtectionLevel(ctx context.Context, state SettingSwitcher, level GlacierVoltageProtectionLevel)(*CmdSetResponse, error)
```
