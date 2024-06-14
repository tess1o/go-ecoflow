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

	//creating new client.
	//Http client can be customized if required (see ecoflow.NewEcoflowClientWithHttpClient(accessKey, secretKey, httpClient))
	client := ecoflow.NewEcoflowClient(accessKey, secretKey)

	//get all linked ecoflow devices
	devices, err := client.GetDeviceList(context.Background())
	if err != nil {
		slog.Error("Cannot get device list", "error", err, "response", devices)
		return
	}

	//for each device get all parameters
	for _, d := range devices.Devices {
		slog.Info("Linked device", "SN", d.SN, "is online", d.Online)

		// get device quote. The list of parameters might be incomplete, raise an Issue if something is missing
		quote, quoteErr := client.GetDeviceAllQuote(context.TODO(), d.SN)
		if quoteErr != nil {
			slog.Error("Cannot get quote for device", "sn", d.SN, "error", quoteErr)
		}
		slog.Info("Quote parameters", "sn", d.SN, "params", quote)

		// get all the parameters as map[string]interface{}. Some values are float64, some are ints and some are []int
		params, paramErr := client.GetDeviceQuoteRawParameters(context.TODO(), d.SN)
		if paramErr != nil {
			slog.Error("Cannot get raw parameters for device", "sn", d.SN, "error", paramErr)
		}
		slog.Info("Raw parameters", "sn", d.SN, "params", params)
	}
}
