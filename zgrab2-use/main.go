package main

import (
	"fmt"
	"github.com/zmap/zgrab2"
	"github.com/zmap/zgrab2/modules/http"
	"net"
)

func main() {
	module := &http.Module{}
	scanner := module.NewScanner()
	// error
	u := uint(443)
	_ = scanner.Init(module.NewFlags().(*http.Flags))
	status, results, err := scanner.Scan(zgrab2.ScanTarget{
		// IP, domain, tag, port
		IP:   net.IP("146.56.155.182"),
		Port: &u,
	})
	fmt.Println(status, results, err)
}
