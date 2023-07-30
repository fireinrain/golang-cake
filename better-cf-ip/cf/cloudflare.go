package cf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

type CloudflareDNS struct{}

var CloudflareConfigValue CloudflareConfig

func init() {
	config := NewCloudflareConfig()
	CloudflareConfigValue = *config
}

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
			config.DnsRecordType = strings.TrimSpace(keyValue[1])
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

type CloudflareDNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type CloudflareDNSResponse struct {
	Result  []CloudflareDNSRecord `json:"result"`
	Success bool                  `json:"success"`
}

// GetAllDNSRecords
//
//	@Description: 获取zoneid 所有dns records
//	@receiver receiver
//	@param DNSType
//	@return []CloudflareDNSRecord
func (receiver *CloudflareDNS) GetAllDNSRecords(DNSType string) []CloudflareDNSRecord {
	// Prepare the API URL
	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", CloudflareConfigValue.ZoneID)

	// Create a new HTTP/2 request
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal("Error creating API request:", err)
	}

	// Set the necessary headers for authentication and content type
	request.Header.Set("X-Auth-Key", fmt.Sprintf("%s", CloudflareConfigValue.ApiKey))
	request.Header.Set("X-Auth-Email", fmt.Sprintf("%s", CloudflareConfigValue.Email))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Make the HTTP request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error making API request:", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading API response:", err)
	}

	// Parse the response JSON into a CloudflareDNSResponse struct
	var cloudflareResponse CloudflareDNSResponse
	if err := json.Unmarshal(body, &cloudflareResponse); err != nil {
		log.Fatal("Error parsing API response:", err)
	}

	// Check if the API request was successful
	if !cloudflareResponse.Success {
		log.Fatal("API request was not successful")
	}

	var results []CloudflareDNSRecord
	// Print the DNS records
	for _, record := range cloudflareResponse.Result {
		//fmt.Printf("ID: %s, Name: %s, Type: %s, Content: %s\n", record.ID, record.Name, record.Type, record.Content)
		if record.Type == DNSType {
			results = append(results, record)
		}
	}
	return results
}

// CheckIfIPAlive
//
//	@Description: 检测cf ip是否可用
//	@receiver receiver
//	@param ipStr
//	@return bool
func (receiver *CloudflareDNS) CheckIfIPAlive(ipStr string) bool {

	return false
}
