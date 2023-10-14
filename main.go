package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/miekg/dns"
)

func GetDNSIPs(domain string) (string, string, error) {
	ipv4Addrs, err := resolveDNS(domain, dns.TypeA)
	if err != nil {
		return "", "", err
	}

	ipv6Addrs, err := resolveDNS(domain, dns.TypeAAAA)
	if err != nil {
		return "", "", err
	}

	return ipv4Addrs[0], ipv6Addrs[0], nil
}

func resolveDNS(domain string, qType uint16) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(domain+".", qType)

	r, _, err := c.Exchange(m, "8.8.8.8:53") // Use a DNS server of your choice
	if err != nil {
		return nil, err
	}

	var addresses []string

	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			addresses = append(addresses, a.A.String())
		} else if aaaa, ok := ans.(*dns.AAAA); ok {
			addresses = append(addresses, aaaa.AAAA.String())
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("No addresses found for %s", domain)
	}

	return addresses, nil
}

func getCurrentIPv4() (string, error) {
	resp, err := http.Get("https://checkipv4.dedyn.io")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getCurrentIPv6() (string, error) {
	resp, err := http.Get("https://checkipv6.dedyn.io")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getToUpdateIPs(domain string) (string, string, error) {
	ipv4, err := getCurrentIPv4()
	if err != nil {
		return "", "", err
	}

	ipv6, err := getCurrentIPv6()
	if err != nil {
		return "", "", err
	}

	curIPv4, curIPv6, err := GetDNSIPs(domain)
	if err != nil {
		return "", "", fmt.Errorf("Error fetching DNS IP addresses: %v", err)
	}

	if ipv4 == curIPv4 && ipv6 == curIPv6 {
		return "", "", fmt.Errorf("No IP address change detected")
	}

	return ipv4, ipv6, nil
}

func updateIPs(domain, token, ipv4, ipv6 string) (string, error) {
	url := "https://update.dedyn.io/"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(domain, token)

	// Add the new IP addresses as query parameters
	query := req.URL.Query()
	query.Add("myipv4", ipv4)
	query.Add("myipv6", ipv6)
	req.URL.RawQuery = query.Encode()

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func dyndns(domain, token string) {
	ipv4, ipv6, err := getToUpdateIPs(domain)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Perform the basic auth HTTP call
	_, err = updateIPs(domain, token, ipv4, ipv6)
	if err != nil {
		fmt.Printf("Error performing HTTP call: %v\n", err)
		return
	}

	fmt.Printf("Successful updated IPs: %s, %s\n", ipv4, ipv6)
	fmt.Println("Waiting for next update in 5 minutes...")
}

func main() {
	domain := os.Getenv("DYNDNS_DOMAIN")
	token := os.Getenv("DYNDNS_TOKEN")
	dyndns(domain, token)

	interval := 5 * time.Minute

	ticker := time.NewTicker(interval)

	tickerChan := ticker.C

	for {
		select {
		case <-tickerChan:
			dyndns(domain, token)
		}
	}
}
