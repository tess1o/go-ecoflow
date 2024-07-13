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

	ctx := context.Background()

	// get set / get functions for power stations. PRO version is not currently implemented
	ps := client.GetPowerStation("SN_HERE")

	//set functions
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

	//get functions
	ps.GetAllParameters(ctx)
	ps.GetParameter(ctx, []string{"mppt.acStandbyMins", "mppt.dcChgCurrent"})

	// get SmartPlug instance with set/get functions
	plug := client.GetSmartPlug("SN_HERE")

	//set functions
	plug.SetRelaySwitch(ctx, ecoflow.SettingEnabled)
	plug.SetIndicatorBrightness(ctx, 1000)
	plug.DeleteScheduledTasks(ctx, 1)

	//get functions
	plug.GetAllParameters(ctx)
	plug.GetParameter(ctx, []string{"2_1.switchSta", "2_1.brightness"})
}
