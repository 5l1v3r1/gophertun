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

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	gophertun "../.."
)

func parseCIDR(cidr string) *net.IPNet {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return &net.IPNet{
		IP:   ip,
		Mask: ipnet.Mask,
	}
}

func main() {
	c := &gophertun.TunTapConfig{
		AllowNameSuffix:       true,
		PreferredNativeFormat: gophertun.FormatEthernet,
	}
	if len(os.Args) >= 2 {
		c.NameHint = os.Args[1]
	}
	t, err := c.Create()
	if err != nil {
		log.Fatalln(err)
	}
	defer t.Close()
	name, err := t.Name()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Device: %s\n", name)
	err = t.SetMTU(65521)
	if err != nil {
		log.Fatalln(err)
	}
	mtu, err := t.MTU()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("MTU:    %d\n", mtu)
	addresses := []*gophertun.IPAddress{
		&gophertun.IPAddress{
			Net:  parseCIDR("10.42.0.2/24"),
			Peer: parseCIDR("10.42.0.3/24"),
		},
		&gophertun.IPAddress{
			Net:  parseCIDR("fd42::2/64"),
			Peer: parseCIDR("fd42::3/64"),
		},
	}
	n, err := t.AddIPAddresses(addresses)
	if err != nil {
		log.Fatalf("Error setting addr[%d]: %s\n", n, err)
	}
	err = t.Open(gophertun.FormatEthernet)
	for {
		p, err := t.Read()
		if err != nil {
			log.Fatalln(err)
		}
		if p == nil {
			break
		}
		fmt.Printf("EtherType: %04x Payload: %x Extra: %x\n", p.EtherType, p.Payload, p.Extra)
	}
}
