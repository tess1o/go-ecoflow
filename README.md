Partial implementation of Ecoflow Rest API that allows to get list of devices and their parameters (quotas).

Link to official documentation: https://developer-eu.ecoflow.com/us/document/introduction
Response mapping: [here](docs/fields_mapping.md.md)
Usage example:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/tess1o/go-ecoflow"
)

const (
	accessKey = "your_access_key"
	secretKey = "your_secret_key"
)

func main() {
	client := ecoflow.NewEcoflowClient(accessKey, secretKey, nil)
	devices, err := client.GetDeviceList(context.Background())
	if err != nil {
		log.Fatalf("Error: %+v\n", err)
		return
	}
	for _, d := range devices.Devices {
		fmt.Printf("Device SN: %s, Online: %d\n", d.SN, d.Online)
		quote, err := client.GetDeviceAllQuote(context.TODO(), d.SN)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}
		log.Printf("Device parameters: %+v\n", quote)
	}
}

```