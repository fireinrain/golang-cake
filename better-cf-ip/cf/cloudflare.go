package cf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CloudflareConfig struct {
	ApiKey string
	Email  string

	ZoneID string
	// DNS record type you want to fetch, e.g., "A", "CNAME", "MX", etc.
	DnsRecordType string
	DnsRecordName string
}

var CloudflareConfigValue CloudflareConfig

func NewCloudflareConfig() *CloudflareConfig {
	config := &CloudflareConfig{
		ApiKey:        "",
		Email:         "",
		ZoneID:        "",
		DnsRecordType: "",
		DnsRecordName: "",
	}
	file, err := os.Open("cloudflare.secret")
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	scanner := bufio.NewScanner(file)

	// Read line by line until the end of the file
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		keyValue := strings.Split(line, "=")
		switch keyValue[0] {

		case "ApiKey":
			config.ApiKey = strings.TrimSpace(keyValue[1])
		case "Email":
			config.Email = strings.TrimSpace(keyValue[1])
		case "ZoneID":
			config.ZoneID = strings.TrimSpace(keyValue[1])
		case "DnsRecordType":
			config.DnsRecordName = strings.TrimSpace(keyValue[1])
		case "DnsRecordName":
			config.DnsRecordName = strings.TrimSpace(keyValue[1])
		default:
			fmt.Println("can not parse this value:", keyValue[0])
		}
		//fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return config
}

func init() {
	config := NewCloudflareConfig()
	CloudflareConfigValue = *config
}
