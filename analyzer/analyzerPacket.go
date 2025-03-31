package analyzer

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type analyzerPacket struct {
	SrcIP       string
	SrcMAC      string
	DstIP       string
	DstMAC      string
	Application string
	Protocol    string
	Timestamp   time.Time
	Length      int
}

func newPacket(packet gopacket.Packet) *analyzerPacket {
	r := &analyzerPacket{}

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth, _ := ethLayer.(*layers.Ethernet)
		r.SrcMAC = eth.SrcMAC.String()
		r.DstMAC = eth.DstMAC.String()
	}

	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		r.SrcIP = ip.SrcIP.String()
		r.DstIP = ip.DstIP.String()
	}

	r.Protocol = getProtocolName(packet)
	r.Application = getApplicationProtocol(packet)
	r.Timestamp = packet.Metadata().Timestamp
	r.Length = packet.Metadata().CaptureLength

	return r
}

func getApplicationProtocol(packet gopacket.Packet) string {
	switch {
	case packet.Layer(layers.LayerTypeDNS) != nil:
		return "DNS"
	case packet.Layer(layers.LayerTypeTLS) != nil:
		return "TLS/SSL"
	case packet.Layer(layers.LayerTypeDHCPv4) != nil:
		return "DHCP"
	default:
		// Try to guess from ports if no specific layer
		if tcp := packet.Layer(layers.LayerTypeTCP); tcp != nil {
			tcp, _ := tcp.(*layers.TCP)
			switch {
			case tcp.SrcPort == 80 || tcp.DstPort == 80:
				return "HTTP (port 80)"
			case tcp.SrcPort == 443 || tcp.DstPort == 443:
				return "HTTPS (port 443)"
			case tcp.SrcPort == 22 || tcp.DstPort == 22:
				return "SSH (port 22)"
			}
		}
		return "Unknown Application"
	}
}

func getProtocolName(packet gopacket.Packet) string {
	// Check for different protocol layers
	// HTTP MISSING
	switch {
	case packet.Layer(layers.LayerTypeTCP) != nil:
		return "TCP"
	case packet.Layer(layers.LayerTypeUDP) != nil:
		return "UDP"
	case packet.Layer(layers.LayerTypeICMPv4) != nil:
		return "ICMPv4"
	case packet.Layer(layers.LayerTypeICMPv6) != nil:
		return "ICMPv6"
	case packet.Layer(layers.LayerTypeDNS) != nil:
		return "DNS"
	case packet.Layer(layers.LayerTypeTLS) != nil:
		return "TLS"
	default:
		return "Unknown"
	}
}

func (ap *analyzerPacket) Print() {
	fmt.Println("Packet: ---")
	fmt.Printf("  * IP: %v -> %v\n", ap.SrcIP, ap.DstIP)
	fmt.Printf("  * MAC: %v -> %v\n", ap.SrcMAC, ap.DstMAC)
	fmt.Printf("  * Protocol: %v | %v\n", ap.Protocol, ap.Application)
	fmt.Printf("  * Length: %d \n", ap.Length)
	fmt.Printf("  * Timestamp: %v \n", ap.Timestamp.String())
}
