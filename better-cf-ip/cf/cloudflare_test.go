package cf

import "testing"

func TestCloudflareDNS_GetAllDNSRecords(t *testing.T) {
	receiver := &CloudflareDNS{}
	receiver.GetAllDNSRecords()
}
