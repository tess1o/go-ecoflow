# Ecoflow API implementation via REST API in Go

## Caution

This library is in heavy development and public API can be changed until a stable version is released.

## About

Partial implementation of Ecoflow Rest API that allows to get list of devices, parameters and set settings.\
The library was tested on Ecoflow Delta 2 and Ecoflow River 2 (I don't have other their products)\
The Ecoflow documentation is not complete and sometimes hard to understand, thus some APIs might not work as expected

## Supported devices:

1. Power Stations (regular ecoflow power stations, like Delta 2, River 2, etc)
2. Power Stations (PRO)
3. Smart Plug
4. PowerStream Micro Inverter
5. Smart Home Panel
6. Wave Air Conditioner

## Features

The library allows to:

1. Get list of all linked devices
2. Get specified parameters from a device
3. Get all parameters from a device
4. Change device's settings

Basically that's all documented features the Ecoflow REST API provides

## Documentation

Link to official documentation: https://developer-eu.ecoflow.com/us/document/introduction

## Usage example

Usage example (also see examples in `examples` folder)

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

	//set functions, see details in documentation to each function.
	ps.SetDcSwitch(ctx, ecoflow.SettingEnabled)
	ps.Set12VDcChargingCurrent(ctx, 100)
	ps.SetAcChargingSettings(ctx, 500, 0)
	ps.SetAcStandByTime(ctx, 60)
	ps.SetBuzzerSilentMode(ctx, ecoflow.SettingDisabled)
	ps.SetCarChargerSwitch(ctx, ecoflow.SettingEnabled)
	ps.SetMaxChargeSoC(ctx, 99)
	ps.SetMinDischargeSoC(ctx, 1)
	ps.SetSoCToTurnOnSmartGenerator(ctx, 50)
	ps.SetSoCToTurnOffSmartGenerator(ctx, 99)
	ps.SetStandByTime(ctx, 60)
	ps.SetCarStandByTime(ctx, 60)
	ps.SetPrioritizePolarCharging(ctx, ecoflow.SettingEnabled)

	//get parameters functions
	ps.GetAllParameters(ctx)
	ps.GetParameter(ctx, []string{"mppt.acStandbyMins", "mppt.dcChgCurrent"})

	// get SmartPlug instance with set/get functions
	plug := client.GetSmartPlug("DEVICE_SERIAL_NUMBER")

	//set functions
	plug.SetRelaySwitch(ctx, ecoflow.SettingEnabled)
	plug.SetIndicatorBrightness(ctx, 1000)
	plug.DeleteScheduledTasks(ctx, 1)

	//get functions
	plug.GetAllParameters(ctx)
	plug.GetParameter(ctx, []string{"2_1.switchSta", "2_1.brightness"})
}

```

## ToDo

I don't have those devices, however I will try to implement the APIs

1. Set API for GLACIER
2. Set API for Power Kits