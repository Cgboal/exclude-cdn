package main

import (
	"bufio"
	"fmt"
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
)

func isURL(candidate string) bool {
	return strings.Contains(candidate, "://")
}

func extractHost(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return u.Host
	}
	return host
}

func CDNFilter() func(string) bool {
	client, err := cdncheck.NewWithCache()
	if err != nil {
		log.Fatal(err)
	}
	resolveName := Resolver()

	return func(line string) bool {
		host := line
		if isURL(line) {
			host = extractHost(line)
		}

		ip := net.ParseIP(host)
		ips := []net.IP{}
		if ip != nil {
			ips = append(ips, ip)
		} else {
			ips = append(ips, resolveName(host)...)
		}
		for _, ip := range ips {
			found, _, err := client.Check(ip)
			if found && err == nil {
				return true
			}
		}
		return false
	}
}

func Resolver() func(string) []net.IP {
	resolver, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		log.Fatal(err)
	}

	return func(name string) []net.IP {
		validIPs := []net.IP{}
		ips, err := resolver.Lookup(name)
		if err != nil {
			return validIPs
		}
		for _, ip := range ips {
			parsedIP := net.ParseIP(ip)
			if parsedIP.To4() == nil {
				continue
			}
			validIPs = append(validIPs, parsedIP)
		}
		return validIPs
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	filter := CDNFilter()

	var wg sync.WaitGroup
	lines := make(chan string)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for line := range lines {
				if !filter(line) {
					fmt.Println(line)
				}
			}
		}()
	}

	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)
	wg.Wait()
}
