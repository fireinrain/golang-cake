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
			fmt.Printf("当前代理ip: %s可用\n", record.Content)
		} else {
			fmt.Printf("当前代理ip: %s不可用\n", record.Content)
		}
	}
}

func TestCloudflareDNS_PatchDNSRecord(t *testing.T) {
	receiver := &CloudflareDNS{}
	receiver.PatchDNSRecord("e2c253d990cb65c327d03a5c03d1ed65", "tw-hnt.ioerror.eu.org", "8.210.117.18")
}
