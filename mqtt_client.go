package ecoflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"io"
	"net/http"
)

const (
	ecoflowLoginUrl         = "https://api.ecoflow.com/auth/login"
	ecoflowScene            = "IOT_APP"
	ecoflowUserType         = "ECOFLOW"
	ecoflowCertificationUrl = "https://api.ecoflow.com/iot-auth/app/certification"
)

type MqttClientConfiguration struct {
	Email            string
	Password         string
	OnConnect        mqtt.OnConnectHandler
	OnConnectionLost mqtt.ConnectionLostHandler
	OnReconnect      mqtt.ReconnectHandler
}

type MqttClient struct {
	Client           mqtt.Client
	connectionConfig *MqttConnectionConfig
}

// NewMqttClient creates a new MQTT client using email and password
// The client is created with the given onConnect, onConnectLost, and messageHandler functions.
// onConnect is executed when we connect to the MQTT broker, in this handler we should subscribe to the topics
// onConnectLost is executed when we are disconnected from MQTT broken
// ClientID is always should be "ANDROID_%uuid%_%user_id%
func NewMqttClient(ctx context.Context, config MqttClientConfiguration) (*MqttClient, error) {
	c, err := getMqttCredentials(ctx, config.Email, config.Password)
	if err != nil {
		return nil, err
	}
	var protocol = c.Protocol
	var broker = c.Url
	var port = c.Port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", protocol, broker, port))
	opts.SetClientID(fmt.Sprintf("ANDROID_%s_%s", uuid.New(), c.UserId))
	opts.SetUsername(c.CertificateAccount)
	opts.SetPassword(c.CertificatePassword)
	opts.SetConnectRetry(true)
	if config.OnConnect != nil {
		opts.OnConnect = config.OnConnect
	}
	if config.OnConnectionLost != nil {
		opts.OnConnectionLost = config.OnConnectionLost
	}
	if config.OnReconnect != nil {
		opts.OnReconnecting = config.OnReconnect
	}
	return &MqttClient{Client: mqtt.NewClient(opts), connectionConfig: c}, nil
}

// GetMqttCredentials get the MQTT credentials using email and password (the same as you use to log in to your Ecoflow app).
// This method allows to get MQTT connection configuration (username/password, host, port, protocol) and subscribe to topic
// to receive devices parameters.
// We first log in to https://api.ecoflow.com/auth/login to receive UserId and Token
// Then log in to https://api.ecoflow.com/iot-auth/app/certification to receive MQTT connection configuration
func getMqttCredentials(ctx context.Context, email, password string) (*MqttConnectionConfig, error) {
	mqttLoginResponse, err := getLoginResponse(ctx, email, password)
	if err != nil {
		return nil, err
	}

	var params = make(map[string]string)
	params["userId"] = mqttLoginResponse.Data.User.UserId

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	certReq, err := http.NewRequestWithContext(ctx, "GET", ecoflowCertificationUrl, bytes.NewReader(jsonParams))
	if err != nil {
		return nil, err
	}

	certReq.Header.Set("Authorization", "Bearer "+mqttLoginResponse.Data.Token)
	certReq.Header.Add("lang", "en_US")

	client := http.Client{}
	resp, err := client.Do(certReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mqttConn *MqttCredentialsResponse
	err = json.Unmarshal(responseBody, &mqttConn)
	if err != nil {
		return nil, err
	}

	c := &mqttConn.Data
	c.UserId = mqttLoginResponse.Data.User.UserId

	return c, nil
}

// getLoginResponse - log in to https://api.ecoflow.com/auth/login with email/password to get UserId and Token,
// which are later used to obtains MQTT connection params
func getLoginResponse(ctx context.Context, email string, password string) (*MqttLoginResponse, error) {
	var params = make(map[string]string)
	params["email"] = email
	params["password"] = base64.StdEncoding.EncodeToString([]byte(password))
	params["scene"] = ecoflowScene
	params["userType"] = ecoflowUserType
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	loginReq, err := http.NewRequestWithContext(ctx, "POST", ecoflowLoginUrl, bytes.NewReader(jsonParams))
	if err != nil {
		return nil, err
	}

	loginReq.Header.Add("lang", "en_US")
	loginReq.Header.Add("content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(loginReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mqttLoginResponse *MqttLoginResponse
	err = json.Unmarshal(responseBody, &mqttLoginResponse)
	if err != nil {
		return nil, err
	}

	return mqttLoginResponse, nil
}

func (m *MqttClient) Connect() error {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// SubscribeForParameters Subscribe to topic to get all device parameters
func (m *MqttClient) SubscribeForParameters(deviceSn string, callback mqtt.MessageHandler) error {
	topicParams := fmt.Sprintf("/app/device/property/%s", deviceSn)
	return m.SubscribeToTopics([]string{topicParams}, callback)
}

// SubscribeToTopics Subscribe to topics
// Assuming that the MQTT client is already connected to the broker
func (m *MqttClient) SubscribeToTopics(topics []string, callback mqtt.MessageHandler) error {
	topicsMap := make(map[string]byte, len(topics))

	for _, t := range topics {
		topicsMap[t] = 1
	}

	token := m.Client.SubscribeMultiple(topicsMap, callback)
	token.Wait()
	return nil
}
