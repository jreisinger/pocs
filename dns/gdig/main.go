// Gdig is similar to dig(1). Stolen from:
// https://github.com/kubernetes-up-and-running/kuard/blob/master/pkg/dnsapi/api.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	t := flag.String("t", "A", "query type")
	flag.Parse()

	name := "."
	if len(flag.Args()) > 0 {
		name = flag.Args()[0]
	}

	resp, err := dnsQuery(*t, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gdig: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(resp)
}

func dnsQuery(t string, name string) (string, error) {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return "", err
	}

	qtype, ok := dns.StringToType[strings.ToUpper(t)]
	if !ok {
		return "", fmt.Errorf("unknown DNS type: %s", t)
	}

	if len(name) == 0 {
		name = "."
	}

	var names []string
	if dns.IsFqdn(name) {
		names = append(names, name)
	} else {
		for _, s := range config.Search {
			names = append(names, name+"."+s)
		}
		names = append(names, name)
	}

	m := new(dns.Msg)
	c := new(dns.Client)

	var r *dns.Msg
	for _, name := range names {
		m.SetQuestion(dns.Fqdn(name), qtype)
		m.RecursionDesired = true
		r, _, err = c.Exchange(m, config.Servers[0]+":"+config.Port)
		if err != nil {
			return "", err
		}
		if len(r.Answer) > 0 {
			return r.String(), nil
		}
	}
	return r.String(), nil
}
