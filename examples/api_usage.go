package main

import (
	"context"
	"fmt"
	"github.com/tess1o/go-ecoflow"
	"log"
	"os"
)

func main() {
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		log.Fatalf("AccessKey and SecretKey are mandatory")
	}

	client := ecoflow.NewEcoflowClient(accessKey, secretKey, nil)
	devices, err := client.GetDeviceList(context.Background())
	if err != nil {
		log.Fatalf("Error: %+v\n", err)
	}
	for _, d := range devices.Devices {
		fmt.Printf("Device SN: %s, Online: %d\n", d.SN, d.Online)
		quote, err := client.GetDeviceAllQuote(context.TODO(), d.SN)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}
		log.Printf("Quote: %+v\n", quote)
	}
}
