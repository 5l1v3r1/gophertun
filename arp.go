/*
  MIT License

  Copyright (c) 2018 Star Brilliant

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  SOFTWARE.
*/

package gophertun

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

func generateMACAddress() net.HardwareAddr {
	result := make(net.HardwareAddr, 6)
	_, _ = rand.Read(result)
	result[0] &= 0xfe
	result[0] |= 0x02
	return result
}

func proxyARP(p *Packet, hwAddr net.HardwareAddr) []*Packet {
	if p.Format != FormatEthernet {
		return nil
	}

	frame := &CodecEthernet{}
	err := frame.Decode(p.Payload)
	if err != nil {
		log.Printf("Warning: %s\n", err)
		return nil
	}

	fmt.Println()
	fmt.Printf("%+v\n", frame)
	fmt.Println()

	switch frame.Type {

	case EtherTypeARP:
		if arp, ok := frame.Payload.(*CodecARP); ok {
			if arp.HardwareType == 1 && arp.HardwareSize == 6 && arp.Opcode == 1 {
				reply := &CodecEthernet{
					Destination: frame.Source,
					Source:      hwAddr,
					Type:        EtherTypeARP,
					Payload: &CodecARP{
						HardwareType:       arp.HardwareType,
						ProtocolType:       arp.ProtocolType,
						HardwareSize:       arp.HardwareSize,
						ProtocolSize:       arp.ProtocolSize,
						Opcode:             2,
						SenderHardwareAddr: hwAddr,
						SenderProtocolAddr: arp.TargetProtocolAddr,
						TargetHardwareAddr: arp.SenderHardwareAddr,
						TargetProtocolAddr: arp.SenderProtocolAddr,
					},
				}
				replyPacket, err := reply.Encode()
				if err != nil {
					log.Printf("Warning: %s\n", err)
					return nil
				}
				return []*Packet{
					&Packet{
						Format:    FormatEthernet,
						EtherType: EtherTypeARP,
						Payload:   replyPacket,
					},
				}
			}
		}

	case EtherTypeIPv6:
		if ipv6, ok := frame.Payload.(*CodecIPv6); ok {
			for layer := ipv6.Payload; layer != nil; layer = layer.NextLayer() {
				if icmpv6, ok := layer.(*CodecICMPv6); ok {
					if icmpv6.Type == 136 {
						// TODO
					}
				}
			}
		}

	}

	return nil
}
