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
        log.Printf("No .env file found: %v\n", err)
    }
}

func logTrafficDetails(protocolCounts map[string]int, bandwidthUsage map[string]int64) {
    logFile, err := os.OpenFile("traffic_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v\n", err)
        return
    }
    defer logFile.Close()

    logger := log.New(logFile, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
    logger.Printf("Protocol counts: %v\n", protocolCounts)
    logger.Printf("Bandwidth usage (bytes): %v\n", bandwidthUsage)
}

func main() {
    loadEnv()

    device := os.Getenv("CAPTURE_DEVICE")
    snapLen := int32(1600)
    promiscuous := false
    devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatalf("Error finding devices: %v\n", err)
    }
    if len(devices) == 0 {
        log.Fatal("No devices found. Ensure you have the necessary permissions.")
    }

    handle, err := pcap.OpenLive(device, snapLen, promiscuous, pcap.BlockForever)
    if err != nil {
        log.Fatalf("Error opening device %s: %v\n", device, err)
    }
    defer handle.Close()

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    packets := packetSource.Packets()

    protocolCounts := make(map[string]int)
    bandwidthUsage := make(map[string]int64)

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    done := make(chan bool)
    go func() {
        for {
            select {
            case packet, ok := <-packets:
                if !ok {
                    log.Println("Packet channel closed.")
                    done <- true
                    return
                }
                processPacket(packet, protocolCounts, bandwidthUsage)
            case <-ticker.C:
                reportAndResetCounts(protocolCounts, bandwidthUsage)
            }
        }
    }()

    <-done
}

func processPacket(packet gopacket.Packet, protocolCounts map[string]int, bandwidthUsage map[string]int64) {
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

func reportAndResetCounts(protocolCounts map[string]int, bandwidthUsage map[string]int64) {
    fmt.Println("Protocol counts in the last 30 seconds: ", protocolCounts)
    fmt.Println("Bandwidth usage in the last 30 seconds (bytes): ", bandwidthUsage)

    logTrafficDetails(protocolCounts, bandwidthUsage)

    for k := range protocolCounts {
        protocolCounts[k] = 0
    }
    for k := range bandwidthUsage {
        bandwidthUsage[k] = 0
    }
}