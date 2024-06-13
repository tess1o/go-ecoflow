package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tess1o/go-ecoflow"
	"log/slog"
	"net/http"
	"os"
	"time"
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
	}

	// configure prometheus scrap interval and metric prefix
	config := ecoflow.PrometheusConfig{Interval: time.Second * 10, Prefix: "ecoflow"}
	client.RecordPrometheusMetrics(&config)

	// start server with metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
