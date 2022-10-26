package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	EndNode Node
)

type Node struct {
	ID           int `json:"id"`
	EndDeviceIds struct {
		DeviceID       string `json:"device_id"`
		ApplicationIds struct {
			ApplicationID string `json:"application_id"`
		} `json:"application_ids"`
		DevEui  string `json:"dev_eui"`
		JoinEui string `json:"join_eui"`
	} `json:"end_device_ids"`
	CorrelationIds []string  `json:"correlation_ids"`
	ReceivedAt     time.Time `json:"received_at"`
	UplinkMessage  struct {
		FPort      int    `json:"f_port"`
		FrmPayload string `json:"frm_payload"`
		RxMetadata []struct {
			GatewayIds struct {
				GatewayID string `json:"gateway_id"`
			} `json:"gateway_ids"`
			Rssi        int     `json:"rssi"`
			ChannelRssi int     `json:"channel_rssi"`
			Snr         float64 `json:"snr"`
		} `json:"rx_metadata"`
		Settings struct {
			DataRate struct {
				Lora struct {
					Bandwidth       int `json:"bandwidth"`
					SpreadingFactor int `json:"spreading_factor"`
				} `json:"lora"`
			} `json:"data_rate"`
		} `json:"settings"`
	} `json:"uplink_message"`
	Simulated bool `json:"simulated"`
}

func main() {
	// t.Run("TestPostService", func(t *testing.T) {
	//try posting json from via http to /webhooks/post
	file := "test.json"
	//get file as *os.Reader
	reader, err := os.Open(file)
	if err != nil {
		// t.Errorf("Error in Open %v", err)
		println("Error in Open %v", err)
	}

	//post json to https://bestiaever.ml/webhooks/post
	//if response is not 200, then the test fails
	request, _ := http.NewRequest(http.MethodPost, "https://bestiaever.ml/webhooks/post", reader)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))
}
