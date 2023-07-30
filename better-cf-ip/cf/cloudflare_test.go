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

func TestCloudflareDNS_CheckIfIPAlive(t *testing.T) {
	const sni = "www.cloudflare.com"
	receiver := &CloudflareDNS{}
	records := receiver.GetAllDNSRecords("A")
	for _, record := range records {
		fmt.Printf("ID: %s, Name: %s, Type: %s, Content: %s\n", record.ID, record.Name, record.Type, record.Content)
		alive, err := receiver.CheckIfIPAlive(record.Content, sni)
		if err != nil {
			fmt.Println("当前ip检测出现错误：", err.Error())
		}
		if alive {
			fmt.Println("当前ip可用")
		} else {
			fmt.Println("当前ip不可用")
		}
	}
}
