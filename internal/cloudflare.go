package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
)

const baseURL = "https://api.cloudflare.com/client/v4"

type Record struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
}

type recordResponse struct {
	Result []Record `json:"result"`
}

// ListDNSRecords fetches DNS records for a given type and host.
func ListDNSRecords(token, zoneID, host string, recordType string) ([]Record, error) {
	requestURL := fmt.Sprintf("%s/zones/%s/dns_records?type=%s&name=%s", baseURL, zoneID, recordType, host)
	body, err := HTTPGet(requestURL, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Accept":        "application/json",
	})
	if err != nil {
		return nil, err
	}

	var resp recordResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// UpdateRecord finds the DNS record of the given type for host and updates its content when needed.
func UpdateRecord(token, zoneID, host, ip string, recordType string) error {
	records, err := ListDNSRecords(token, zoneID, host, recordType)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("host not found")
	}
	rec := records[0]
	if rec.Content == ip {
		// No-op when content matches desired IP
		return nil
	}

	updatePayload := struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		Proxied bool   `json:"proxied"`
	}{Type: rec.Type, Name: rec.Name, Content: ip, Proxied: rec.Proxied}

	payloadBytes, err := json.Marshal(updatePayload)
	if err != nil {
		return err
	}
	_, err = HTTPPut(
		fmt.Sprintf("%s/zones/%s/dns_records/%s", baseURL, zoneID, rec.ID),
		payloadBytes,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
			"Accept":        "application/json",
			"Content-Type":  "application/json",
		},
	)
	return err
}
