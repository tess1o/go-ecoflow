Partial implementation of Ecoflow Rest API that allows to get list of devices and their parameters (quotas).

Link to official documentation: https://developer-eu.ecoflow.com/us/document/introduction \
Response mapping: [here](docs/fields_mapping.md)
No external dependencies.
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

	//creating new client. Http client can be customized if required
	client := ecoflow.NewEcoflowClient(accessKey, secretKey)

	//get all linked ecoflow devices
	devices, err := client.GetDeviceList(context.Background())
	if err != nil {
		slog.Error("Cannot get device list", "error", err)
		return
	}

	//for each device get all parameters
	for _, d := range devices.Devices {
		slog.Info("Linked device", "SN", d.SN, "is online", d.Online)

		quote, quoteErr := client.GetDeviceAllQuote(context.TODO(), d.SN)
		if quoteErr != nil {
			slog.Error("Cannot get quote for device", "sn", d.SN, "error", quoteErr)
		}
		slog.Info("Quote parameters", "sn", d.SN, "params", quote)

		params, paramErr := client.GetDeviceQuoteRawParameters(context.TODO(), d.SN)
		if paramErr != nil {
			slog.Error("Cannot get raw parameters for device", "sn", d.SN, "error", paramErr)
		}
		slog.Info("Raw parameters", "sn", d.SN, "params", params)
	}
}

```