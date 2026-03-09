package sniffer

import (
	"fmt"
	"log"

	"albion/common/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// RealSniffer captures actual network traffic.
type RealSniffer struct {
	device  string
	updates chan models.MarketOrder
	stop    chan struct{}
	handle  *pcap.Handle
}

// NewRealSniffer creates a new RealSniffer.
// If device is empty, it attempts to find the first suitable device.
func NewRealSniffer(device string) (*RealSniffer, error) {
	if device == "" {
		// Auto-detect
		devs, err := pcap.FindAllDevs()
		if err != nil {
			return nil, err
		}
		if len(devs) == 0 {
			return nil, fmt.Errorf("no network devices found")
		}
		// Pick the first one with an IP address (naive)
		for _, d := range devs {
			if len(d.Addresses) > 0 {
				device = d.Name
				log.Printf("Auto-selected device: %s (%s)", d.Name, d.Description)
				break
			}
		}
	}

	return &RealSniffer{
		device:  device,
		updates: make(chan models.MarketOrder, 100),
		stop:    make(chan struct{}),
	}, nil
}

func (s *RealSniffer) Start() error {
	log.Printf("Starting Real Sniffer on device: %s", s.device)

	handle, err := pcap.OpenLive(s.device, 1600, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	s.handle = handle

	// Filter for Albion's Photon UDP traffic (usually port 5056)
	// We might need to listen to both directions.
	if err := handle.SetBPFFilter("udp port 5056"); err != nil {
		return fmt.Errorf("failed to set BPF filter: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	
	go func() {
		for {
			select {
			case <-s.stop:
				handle.Close()
				return
			case packet := <-packetSource.Packets():
				s.processPacket(packet)
			}
		}
	}()

	return nil
}

func (s *RealSniffer) Stop() {
	close(s.stop)
	close(s.updates)
}

func (s *RealSniffer) Updates() <-chan models.MarketOrder {
	return s.updates
}

func (s *RealSniffer) processPacket(packet gopacket.Packet) {
	app := packet.ApplicationLayer()
	if app != nil {
		payload := app.Payload()
		// Logic: If we see a payload > 10 bytes on UDP 5056, it's likely Albion.
		if len(payload) > 10 {
			// Emit a "Signal" order to show connectivity
			// We limit this to once every few seconds in a real app to avoid spam,
			// but for now let's just emit one to prove it works.
			
			// We can't really "fake" a market order without confusing the user.
			// But we can log it.
			log.Printf("Captured Albion Packet: %d bytes", len(payload))
			
			// Optional: Emit a special "System" order?
			// s.updates <- models.MarketOrder{
			// 	ItemID: "SYSTEM_ALBION_TRAFFIC_DETECTED",
			// 	UnitPrice: len(payload),
			// 	Source: "real_sniffer",
			// 	Confidence: "system",
			// }
		}
	}
}
