package ecoflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// PrometheusConfig represents the configuration for recording Prometheus metrics.
type PrometheusConfig struct {

	// Prefix represents the metric prefix to be used when recording Prometheus metrics. It's a good idea to use value "ecoflow"
	Prefix string

	// Interval represents scrap interval at which Prometheus metrics are recorded.
	Interval time.Duration
}

// RecordPrometheusMetrics is responsible for recording Prometheus metrics at a specified interval.
// It creates a ticker with the provided interval and starts a goroutine that iterates through the list of devices,
// retrieves their quotas, and handles each metric using the handleOneMetric and handleMetrics functions.
// If any error occurs while getting the device list or quota, an error is logged and the process continues with the next device.
// If the metric name generation fails, an error is logged and the metric is skipped.
// If the metric value is an array, it is skipped.
// If the metric value is a float64, it is set as the gauge value.
// If the metric value cannot be converted to a float64, an error is logged and the metric is skipped.
//
// Args:
//
//	config: The PrometheusConfig object.
//
// Returns: None.
func (c *Client) RecordPrometheusMetrics(config *PrometheusConfig) {
	ticker := time.NewTicker(config.Interval)
	var metrics = make(map[string]prometheus.Gauge)
	go func() {
		for _ = time.Now(); ; _ = <-ticker.C {
			devices, err := c.GetDeviceList(context.Background())
			if err != nil || devices.Code != "0" {
				slog.Error("Cannot get devices list", "error", err)
				continue
			}

			for _, dev := range devices.Devices {
				quota, quoteErr := c.GetDeviceAllQuote(context.Background(), dev.SN)
				if quoteErr != nil || quota.Code != "0" {
					slog.Error("Cannot get device quota", "device", dev.SN, "error", quoteErr)
					continue
				}

				handleOneMetric(metrics, "online", float64(dev.Online), config, dev)

				handleMetrics(metrics, quota.Data.Pd, config, dev)
				handleMetrics(metrics, quota.Data.Inv, config, dev)
				handleMetrics(metrics, quota.Data.Mppt, config, dev)
				handleMetrics(metrics, quota.Data.BmsEmsStatus, config, dev)
				handleMetrics(metrics, quota.Data.BmsBmsStatus, config, dev)
				handleMetrics(metrics, quota.Data.BmsSlave, config, dev)
			}
		}
	}()
}

// handleMetrics iterates through the parameterGroup map, calling handleOneMetric for each key-value pair.
// handleOneMetric is responsible for handling a single metric by generating a metric name,
// creating a new gauge if necessary, and setting the gauge value.
// If the metric name generation fails, an error is logged and the function returns.
// If the metric already exists, the gauge value is updated.
// If the metric value is an array, it is skipped.
// If the metric value is a float64, it is set as the gauge value.
// If the metric value cannot be converted to a float64, an error is logged and the metric is skipped.
//
// Args:
//
//	metrics: A map of device metric names to prometheus Gauge objects.
//	parameterGroup: The parameterGroup map containing the metrics to be handled.
//	config: The PrometheusConfig object.
//	dev: The DeviceInfo object.
//
// Returns: None.
func handleMetrics(metrics map[string]prometheus.Gauge, parameterGroup any, config *PrometheusConfig, dev DeviceInfo) {
	if parameterGroup == nil {
		slog.Debug("No parameters provided")
		return
	}

	var params map[string]interface{}
	paramBytes, err := json.Marshal(parameterGroup)
	if err != nil {
		slog.Error("Unable to marshal parameters", "error", err)
		return
	}
	err = json.Unmarshal(paramBytes, &params)
	if err != nil {
		slog.Error("Unable to unmarshal parameters", "error", err)
		return
	}

	for field, val := range params {
		handleOneMetric(metrics, field, val, config, dev)
	}
}

// handleOneMetric handles a single metric by generating a metric name,
// creating a new gauge if necessary, and setting the gauge value.
// If the metric name generation fails, an error is logged and the function returns.
// If the metric already exists, the gauge value is updated.
// If the metric value is an array, it is skipped.
// If the metric value is a float64, it is set as the gauge value.
// If the metric value cannot be converted to a float64, an error is logged and the metric is skipped.
//
// Args:
//
//	metrics: A map of device metric names to prometheus Gauge objects.
//	field: The field name of the metric.
//	val: The metric value.
//	config: The PrometheusConfig object.
//	dev: The DeviceInfo object.
//
// Returns: None.
func handleOneMetric(metrics map[string]prometheus.Gauge, field string, val interface{}, config *PrometheusConfig, dev DeviceInfo) {
	metricName, deviceMetricName, err := generateMetricName(field, config.Prefix, dev.SN)
	if err != nil {
		slog.Error("Unable to generate metric name", "metric", field)
		return
	}
	gauge, ok := metrics[deviceMetricName]
	if !ok {
		slog.Debug("Adding new metric", "metric", metricName, "device", dev.SN)
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metricName,
			ConstLabels: map[string]string{
				"device": dev.SN,
			},
		})
		prometheus.MustRegister(gauge)
		metrics[deviceMetricName] = gauge
	} else {
		slog.Debug("Updating metric", "metric", metricName, "value", val, "device", dev.SN)
	}
	_, ok = val.([]interface{})
	if ok {
		slog.Debug("The value is an array, skipping it", "metric", metricName)
		return
	}
	floatVal, ok := val.(float64)
	if ok {
		gauge.Set(floatVal)
	} else {
		slog.Error("Unable to convert value to float, skipping metric", "value", val, "metric", metricName)
	}
}

// generateMetricName takes a rawMetric string (for example "pd.wireUsedTime") received from Ecoflow Rest API.
// the function returns two values: metricName and deviceMetric name.
// metricName is used as a metric name in prometheus
// deviceMetricName is used as unique key from all metrics (device serial number + metric name). It's required to store all metrics in a map.
// The function converts rawMetric to prometheus compatible value (pd_wire_used_time)
// As a result it returns two values:
// metricName: prefix + prometheus compatible value (ecoflow_pd_wire_used_time)
// deviceMetricName: device serial number + metricName (R13124123123213_ecoflow_pd_wire_used_time)
// Example:
// rawMetric: pd.wireUsedTime
// prefix: ecoflow
// deviceSn: R13124123123213
//
// return values:
// metricName: ecoflow_pd_wire_used_time
// deviceMetricName: R13124123123213_ecoflow_pd_wire_used_time

func generateMetricName(rawMetric string, prefix string, deviceSn string) (string, string, error) {
	prometheusName, err := convertEcoflowKeyToPrometheusName(rawMetric)
	if err != nil {
		return "", "", err
	}
	metricName := prefix + "_" + prometheusName
	deviceMetricName := deviceSn + "_" + metricName
	return metricName, deviceMetricName, nil
}

// convertEcoflowKeyToPrometheusName takes a metricKey string and converts it to a Prometheus compatible name.
// It replaces the dot "." with an underscore "_" and converts any uppercase letters to lowercase,
// adding an underscore before each uppercase letter except when it's preceded by an underscore.
// The converted metricKey string must adhere to the Prometheus data model pattern [a-zA-Z_:][a-zA-Z0-9_:]*.
// If the conversion fails, an error is returned.
//
// Example:
//
//	metricKey: "pd.wireUsedTime"
//	converted: "pd_wire_used_time"
//
// Args:
//
//	metricKey: The original metric key string.
//
// Returns:
//
//	The converted metric key string, nil if successful.
//	Returns an error if the conversion fails.
func convertEcoflowKeyToPrometheusName(metricKey string) (string, error) {
	key := strings.Replace(metricKey, ".", "_", -1)
	runes := []rune(key)
	var newKey bytes.Buffer

	newKey.WriteString(strings.ToLower(string(runes[0])))

	for _, character := range runes[1:] {
		if unicode.IsUpper(character) && !strings.HasSuffix(newKey.String(), "_") {
			newKey.WriteString("_")
		}
		newKey.WriteString(strings.ToLower(string(character)))
	}

	matches, err := regexp.MatchString("[a-zA-Z_:][a-zA-Z0-9_:]*", newKey.String())
	if err != nil {
		return "", err
	}

	if !matches {
		return "", fmt.Errorf("cannot convert payload key %s to comply with the Prometheus data model", metricKey)
	}

	return newKey.String(), nil
}
