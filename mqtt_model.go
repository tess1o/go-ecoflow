package ecoflow

// MqttConnectionConfig represents the configuration for MQTT connection.
// It contains the following fields:
// - CertificateAccount: the account for the MQTT certificate
// - CertificatePassword: the password for the MQTT certificate
// - Url: the URL of the MQTT broker
// - Port: the port number for the MQTT connection
// - Protocol: the protocol for the MQTT connection
// - UserId: the user ID received from login call
type MqttConnectionConfig struct {
	CertificateAccount  string `json:"certificateAccount"`
	CertificatePassword string `json:"certificatePassword"`
	Url                 string `json:"url"`
	Port                string `json:"port"`
	Protocol            string `json:"protocol"`
	UserId              string `json:"userId"`
}

// MqttCredentialsResponse represents the response structure for MQTT credentials.
// It contains the following fields:
// - Code: the code of the response. 0: success, otherwise - error code.
// - Message: the message returned in the response
// - Data: the MQTT connection configuration
type MqttCredentialsResponse struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Data    MqttConnectionConfig `json:"data"`
}

// MqttDeviceParams represents the device parameters received from MQTT topic
// Params map is a key/value map where key is parameter name and value is its value
type MqttDeviceParams struct {
	Id         int64                  `json:"id"`
	Timestamp  int                    `json:"timestamp"`
	ModuleType string                 `json:"moduleType"`
	Params     map[string]interface{} `json:"params"`
}

// MqttLoginResponse when we log in to ecoflow rest api via email/password.
// We actually use only `Data.Token` and `Data.User.UserId`
type MqttLoginResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		User struct {
			UserId        string `json:"userId"`
			Email         string `json:"email"`
			Name          string `json:"name"`
			Icon          string `json:"icon"`
			State         int    `json:"state"`
			Regtype       string `json:"regtype"`
			CreateTime    string `json:"createTime"`
			Destroyed     string `json:"destroyed"`
			RegisterLang  string `json:"registerLang"`
			Source        string `json:"source"`
			Administrator bool   `json:"administrator"`
			Appid         int    `json:"appid"`
			CountryCode   string `json:"countryCode"`
		} `json:"user"`
		Token string `json:"token"`
	} `json:"data"`
}
