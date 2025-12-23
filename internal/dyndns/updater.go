package dyndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"codeberg.org/maltech/desec-dyndns-client/internal/config"
	"codeberg.org/maltech/desec-dyndns-client/internal/metrics"
)

func getRRset(cfg *config.Config, recordType string) ([]string, error) {
	url := fmt.Sprintf("https://desec.io/api/v1/domains/%s/rrsets/%s/%s/", cfg.Domain, cfg.Subname, recordType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		// RRset exists
		var rrset struct {
			Records []string `json:"records"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&rrset); err != nil {
			return nil, err
		}
		return rrset.Records, nil

	case 404:
		// No RRset yet — treat as empty
		return []string{}, nil

	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s (status %d)", string(body), resp.StatusCode)
	}
}

func updateRRset(cfg *config.Config, recordType, ip string) error {
	url := fmt.Sprintf("https://desec.io/api/v1/domains/%s/rrsets/", cfg.Domain)

	payload := map[string]interface{}{
		"subname": cfg.Subname,
		"type":    recordType,
		"ttl":     cfg.TTL,
		"records": []string{ip},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token "+strings.TrimSpace(cfg.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("INFO: %s record for %s.%s updated to %s", recordType, cfg.Subname, cfg.Domain, ip)
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("API error: %s (status %d)", string(body), resp.StatusCode)
}

func getPublicIP(version string) (string, error) {
	url := "https://checkipv4.dedyn.io"
	if version == "AAAA" {
		url = "https://checkipv6.dedyn.io"
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	return string(ip), err
}

func (a *App) update(recordType string) {
	ip, err := getPublicIP(recordType)
	if err != nil {
		log.Printf("WARN: failed to get IPv4: %v", err)
	}
	rrset, err := getRRset(a.cfg, recordType)
	if err != nil {
		log.Printf("WARN: failed to fetch %s record: %v", recordType, err)
		metrics.Failure.WithLabelValues(recordType, a.cfg.Hostname).Inc()
		return
	}

	// Update if no rrset or ip not in rrset
	if len(rrset) == 0 || ip != rrset[0] {
		if err := updateRRset(a.cfg, recordType, ip); err != nil {
			metrics.Failure.WithLabelValues(recordType).Inc()
			log.Printf("WARN: failed to update %s: %v", recordType, err)
		}
		metrics.Success.WithLabelValues(recordType, a.cfg.Hostname)
		metrics.Last.WithLabelValues(recordType, a.cfg.Hostname).SetToCurrentTime()
	} else {
		log.Printf("INFO: IP %s of %s record %s.%s unchanged, no update needed", ip, recordType, a.cfg.Subname, a.cfg.Domain)
	}
}
