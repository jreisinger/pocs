package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

var (
	addr   = "localhost:53000"
	domain = "ZG5ZC2VJDXJPDHKK.COM."
	ipaddr = "1.2.3.4"
)

func main() {
	// Create a new DNS server
	server := &dns.Server{Addr: addr, Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)

	// Start the DNS server
	log.Println("starting DNS server at", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start DNS server: %v", err)
	}
}

func handleDNSRequest(res dns.ResponseWriter, req *dns.Msg) {
	// Create a new DNS message
	msg := dns.Msg{}
	msg.SetReply(req)
	msg.Authoritative = true

	// Handle each question in the request
	for _, q := range req.Question {
		switch q.Qtype {
		case dns.TypeA:
			handleARecord(&msg, q)
		}
	}

	// Send the response
	res.WriteMsg(&msg)
}

func handleARecord(msg *dns.Msg, q dns.Question) {
	// Extract and decode data from the query
	parts := strings.SplitN(q.Name, ".", 2)
	if len(parts) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			log.Printf("failed to decode base64 string: %v", err)
		}
		fmt.Println("decoded info:", string(decoded))
	}

	// Answer the query
	if strings.HasSuffix(q.Name, domain) {
		// Create a new A record
		rr := &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    600,
			},
			A: net.ParseIP(ipaddr),
		}
		msg.Answer = append(msg.Answer, rr)
	}
}
