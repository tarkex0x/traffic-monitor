package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	loadEnv()

	device := os.Getenv("CAPTURE_DEVICE")
	snapLen := int32(1600)
	promiscuous := false
	err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	handle, err := pcap.OpenLive(device, snapLen, promiscuous, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()

	protocolCounts := make(map[string]int)
	bandwidthUsage := make(map[string]int64)

	for packet := range packets {
		if networkLayer := packet.NetworkLayer(); networkLayer != nil {
			protocol := networkLayer.NetworkFlow().EndpointType().String()
			protocolCounts[protocol]++

			bandwidthUsage[protocol] += int64(len(packet.Data()))
		}

		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			if string(appLayer.Payload()) == "unusual DNS query pattern" {
				fmt.Println("Alert: Unusual DNS query pattern detected!")
			}
		}
	}

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Protocol counts in the last 30 seconds: ", protocolCounts)
			fmt.Println("Bandwidth usage in the last 30 seconds (bytes): ", bandwidthUsage)

			for k := range protocolCounts {
				protocolCounts[k] = 0
			}
			for k := range bandwidthUsage {
				bandwidthUsage[k] = 0
			}
		}
	}
}