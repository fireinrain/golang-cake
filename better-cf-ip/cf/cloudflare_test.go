package cf

import (
	"fmt"
	"testing"
)

func TestCloudflareDNS_GetAllDNSRecords(t *testing.T) {
	receiver := &CloudflareDNS{}
	records := receiver.GetAllDNSRecords("A")
	for _, record := range records {
		fmt.Printf("ID: %s, Name: %s, Type: %s, Content: %s\n", record.ID, record.Name, record.Type, record.Content)
	}
}
