package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var n = flag.Int("n", 0, "look up only against first `n` nameservers")
var t = flag.String("t", "A", "DNS `type` to look up")

func main() {
	flag.Parse()

	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0) // no timestamp

	if len(flag.Args()) != 1 {
		log.Fatalf("supply one FQDN (e.g. example.com) to look up")
	}

	var fqdn = flag.Arg(0)
	var stats statistics
	var servers nameservers

	servers.add("1.1.1.1", "1.0.0.1") // Cloudflare
	servers.add("8.8.8.8", "8.8.4.4") // Google
	if err := servers.getLocal(); err != nil {
		log.Printf("getting local nameservers: %v\n", err)
	}
	if err := servers.getPublic("https://public-dns.info/nameservers.txt"); err != nil {
		log.Printf("getting public nameservers: %v\n", err)
	}
	servers.dedup()
	if *n != 0 {
		servers = servers[:*n]
	}

	stats.totalServers = len(servers)

	serversCh := make(chan string)
	resultsCh := make(chan statistics)

	t := strings.ToUpper(*t)
	dnsType := dns.StringToType[t]
	if dnsType == 0 {
		log.Fatalf("unsupported DNS type: %s", t)
	}

	// Spin up 100 workers to make lookups.
	for i := 0; i < 100; i++ {
		go worker(fqdn, dnsType, serversCh, resultsCh)
	}

	// Send workers servers to make lookups at.
	go func() {
		for _, s := range servers {
			serversCh <- s
		}
	}()

	for range servers {
		s := <-resultsCh
		stats.failedServers += s.failedServers
		stats.totalResponses += s.totalResponses
		stats.failedResponses += s.failedResponses
		stats.emptyResponses += s.emptyResponses
	}

	close(serversCh)
	close(resultsCh)

	stats.printSummary()
}

// statistics holds statistics about DNS responses.
type statistics struct {
	totalServers    int
	failedServers   int
	totalResponses  int
	failedResponses int
	emptyResponses  int
}

func (s *statistics) printSummary() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	format := "%-24s %4.0f%% (%d/%d)\n"
	fmt.Printf(format, "Failed nameservers",
		s.failedServersPercentage(), s.failedServers, s.totalServers)
	fmt.Printf(format, "Failed responses",
		s.failedResponsesPercentage(), s.failedResponses, s.totalResponses)
	fmt.Printf(format, "Empty responses",
		s.emptyResponsesPercentage(), s.emptyResponses, s.totalResponses-s.failedResponses)
}

func (s *statistics) failedServersPercentage() float64 {
	return float64(s.failedServers) / float64(s.totalServers) * 100
}

func (s *statistics) failedResponsesPercentage() float64 {
	return float64(s.failedResponses) / float64(s.totalResponses) * 100
}

func (s *statistics) emptyResponsesPercentage() float64 {
	return float64(s.emptyResponses) / float64(s.totalResponses-s.failedResponses) * 100
}

// nameservers to make lookups against.
type nameservers []string

func (n *nameservers) add(servers ...string) {
	*n = append(*n, servers...)
}

func (n *nameservers) dedup() {
	orig := *n
	*n = []string{} // empty the slice
	seen := make(map[string]struct{})
	for _, s := range orig {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		*n = append(*n, s)
	}
}

func (n *nameservers) getLocal() error {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return err
	}
	*n = append(*n, config.Servers...)
	return nil
}

func (n *nameservers) getPublic(url string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Select only IPv4 nameservers.
	for _, server := range strings.Split(string(b), "\n") {
		ipv4 := net.ParseIP(server).To4()
		if ipv4 != nil {
			*n = append(*n, ipv4.String())
		}
	}

	return nil
}

func worker(fqdn string, dnsType uint16, servers chan string, results chan statistics) {
	for server := range servers {
		stats := lookup(fqdn, server, dnsType)
		results <- stats
	}
}

func lookup(fqdn, server string, dnsType uint16) statistics {
	var stats statistics
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dnsType)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, net.JoinHostPort(server, "53"))
	if r == nil { // server issues
		stats.failedServers++
		fmt.Printf("lookup at %-15s %v\n", server, err.Error())
		return stats
	}

	stats.totalResponses++

	nRRs := len(r.Answer)
	rCode := dns.RcodeToString[r.Rcode]

	if rCode != "NOERROR" {
		stats.failedResponses++
		fmt.Printf("lookup at %-15s %v\n", server, rCode)
		return stats
	}

	format := "lookup at %-15s %d RR"
	if nRRs == 0 || nRRs > 1 {
		format += "s"
	}
	fmt.Printf(format+"\n", server, nRRs)

	if nRRs < 1 {
		stats.emptyResponses++
	}

	return stats
}
