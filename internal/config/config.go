package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Config struct {
	Domain      string
	Subname     string
	Hostname    string
	Token       string
	Interval    time.Duration
	TTL         int
	MetricsAddr string
}

func splitHost(host string) (subname, domain string, err error) {
	etldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", "", err
	}
	subname = strings.TrimSuffix(host, "."+etldPlusOne)
	if subname == "" {
		subname = "@"
	}
	return subname, etldPlusOne, nil
}

func Load() *Config {
	hostname := os.Getenv("DYNDNS_HOSTNAME")
	token := os.Getenv("DYNDNS_TOKEN")

	// Deprecated variable
	if old := os.Getenv("DYNDNS_DOMAIN"); old != "" {
		if hostname == "" {
			hostname = old
		}
		log.Printf("WARN: DYNDNS_DOMAIN is deprecated. Please use DYNDNS_HOSTNAME instead.")
	}

	if hostname == "" {
		log.Fatal("FATAL: DYNDNS_HOSTNAME is not set")
	}
	if token == "" {
		log.Fatal("FATAL: DYNDNS_TOKEN is not set")
	}

	subname, domain, err := splitHost(hostname)
	if err != nil {
		log.Fatalf("FATAL: failed to split hostname '%s': %v", hostname, err)
	}

	// Defaults
	defaultInterval := 5 * time.Minute
	defaultTTL := 300
	defaultMetricsAddr := ":9333"

	// Read interval from env
	interval := defaultInterval
	if s := os.Getenv("DYNDNS_INTERVAL"); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			interval = d
		} else {
			log.Printf("WARN: invalid DYNDNS_INTERVAL '%s', using default %v", s, defaultInterval)
		}
	}

	// Read TTL from env
	ttl := defaultTTL
	if s := os.Getenv("DYNDNS_TTL"); s != "" {
		if t, err := strconv.Atoi(s); err == nil && t > 0 {
			ttl = t
		} else {
			log.Printf("WARN: invalid DYNDNS_TTL '%s', using default %d", s, defaultTTL)
		}
	}

	// Read MetricsAddr from env
	metricsAddr := defaultMetricsAddr
	if s := os.Getenv("DYNDNS_METRICS_ADDR"); s != "" {
		if !strings.Contains(s, ":") {
			log.Fatalf("FATAL: invalid DYNDNS_METRICS_ADDR '%s'", s)
		}
		metricsAddr = s
	}

	return &Config{
		Domain:      domain,
		Subname:     subname,
		Hostname:    hostname,
		Token:       token,
		Interval:    interval,
		TTL:         ttl,
		MetricsAddr: metricsAddr,
	}
}
