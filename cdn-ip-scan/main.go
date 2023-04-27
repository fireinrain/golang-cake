package main

import (
	"crypto/tls"
	"fmt"
	_ "net"
	"strings"
	"sync"
	"time"
)

// 全球ip信息
// http://ipblock.chacuo.net/

type CheckResult struct {
	Ip        string
	IsProxyIp bool
	Error     error
}

func main() {
	sampleIps := []string{
		"18.163.249.175",
		"128.14.140.254",
		"152.69.204.164",
		"107.172.242.3",
		"23.237.33.106",
		"173.205.94.4",
		"104.223.102.254",
		"128.14.142.176",
		"211.72.35.110",
		"119.36.161.40",
		"45.64.22.53",
		"45.64.22.56",
		"45.64.22.21",
		"45.64.22.22",
		"45.64.22.23",
		"45.64.22.6",
	}
	resultsChan := make(chan CheckResult)
	waitGroup := sync.WaitGroup{}

	for _, ip := range sampleIps {
		go func(ip string) {
			waitGroup.Add(1)
			SNIChecker(ip, "www.cloudflare.com", resultsChan)
			waitGroup.Done()
		}(ip)
	}
	go func() {
		waitGroup.Wait()
		close(resultsChan)
	}()
	//channel没关闭 容易阻塞
	for result := range resultsChan {
		if result.Error != nil {
			fmt.Println("checker error:", result.Error)
		}
		if result.IsProxyIp {
			fmt.Printf("ip: %s is a proxy for cloudflare\n", result.Ip)
		}
	}
}
func SNIChecker(ipStr string, serverName string, resultChan chan CheckResult) {
	// Replace <IP> with the target IP address.
	addr := fmt.Sprintf("%s:443", ipStr)
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: serverName,
	})
	if err != nil {
		fmt.Printf("Error connecting to server: %s to %s\n", err, ipStr)
		resultChan <- CheckResult{
			Ip:        ipStr,
			IsProxyIp: false,
			Error:     err,
		}
		return
		//return false, errors.New("error connecting to server")
	}
	defer conn.Close()

	// Print the server certificate details.
	certs := conn.ConnectionState().PeerCertificates
	//for i, cert := range certs {
	//	fmt.Printf("Certificate %d:\n", i+1)
	//	fmt.Printf("  Subject: %s\n", cert.Subject.CommonName)
	//	fmt.Printf("  Issuer: %s\n", cert.Issuer.CommonName)
	//	fmt.Printf("  Valid from: %s\n", cert.NotBefore)
	//	fmt.Printf("  Valid until: %s\n", cert.NotAfter)
	//	fmt.Println()
	//}

	isPassed := false
	for _, cert := range certs {
		if strings.Contains(cert.Subject.CommonName, serverName) && cert.NotAfter.After(time.Now()) {
			//return true, nil
			isPassed = true
			break
		}

	}
	if isPassed {
		resultChan <- CheckResult{
			Ip:        ipStr,
			IsProxyIp: true,
			Error:     nil,
		}
	} else {
		resultChan <- CheckResult{
			Ip:        ipStr,
			IsProxyIp: false,
			Error:     nil,
		}
	}
}

//IPv4地址空间中有一部分地址是被保留或未分配的，这些地址不能被用于互联网的通信。以下是IPv4地址空间中保留或未分配的地址：
//
//0.0.0.0/8：保留地址，用于表示本地网络。
//10.0.0.0/8：私有地址，用于局域网。
//127.0.0.0/8：保留地址，用于回环测试。
//169.254.0.0/16：自动配置地址，用于本地网络设备自动获取IP地址。
//172.16.0.0/12：私有地址，用于局域网。
//192.0.0.0/24：保留地址，用于IPv4-IPv6转换。
//192.0.2.0/24：保留地址，用于示例和测试网络。
//192.88.99.0/24：保留地址，用于IPv6转换。
//192.168.0.0/16：私有地址，用于局域网。
//198.18.0.0/15：保留地址，用于性能测试。
//198.51.100.0/24：保留地址，用于示例和测试网络。
//203.0.113.0/24：保留地址，用于示例和测试网络。
//224.0.0.0/4：保留地址，用于多播通信。
//240.0.0.0/4：保留地址，用于未分配的地址空间。
//这些地址可以在互联网中自由使用，但它们不能被路由或用于与互联网的通信。任何尝试从这些地址中发送数据包的尝试都将被丢弃或拒绝。

//枚举ipv4 排除保留地址

//枚举ipv6 排除保留地址
