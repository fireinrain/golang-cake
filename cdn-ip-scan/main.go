package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	_ "net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// curl -svo /dev/null -H 'Host: www.cloudflare.com' http://220.130.80.179
//openssl s_client -connect 220.130.80.179:443 -servername www.cloudflare.com

// echo 220.130.80.179 | zgrab2 tls --server-name www.cloudflare.com | grep "cloudflare" | jq '.ip'

// zmap 快速判断443端口开放
// sudo zmap -B 100K -p 443 211.72.0.0/16 -o results.csv

// IpRange
//
//	IpRange
//	@Description: IP range表示
type IpRange struct {
	IPStart string `json:"IPStart,omitempty"`
	IPEnd   string `json:"IPEnd,omitempty"`
	IPCount int    `json:"IPCount,omitempty"`
}

// IpRange2CIDR
//
//	@Description: 将ip range转化为CIDR表示
//	@receiver r
//	@return cidr
//	@return err
func (r IpRange) IpRange2CIDR() (cidr string, err error) {
	startIP := r.IPStart
	endIP := r.IPEnd
	// Parse start and end IPs
	startIPParts := strings.Split(startIP, ".")
	endIPParts := strings.Split(endIP, ".")

	if len(startIPParts) != 4 || len(endIPParts) != 4 {
		return "", fmt.Errorf("invalid IP address format")
	}

	// Convert IP parts to integers
	var startIPInt, endIPInt uint32
	for i := 0; i < 4; i++ {
		startOctet, err := strconv.Atoi(startIPParts[i])
		if err != nil {
			return "", fmt.Errorf("invalid IP address format")
		}
		endOctet, err := strconv.Atoi(endIPParts[i])
		if err != nil {
			return "", fmt.Errorf("invalid IP address format")
		}

		if startOctet < 0 || startOctet > 255 || endOctet < 0 || endOctet > 255 {
			return "", fmt.Errorf("invalid IP address range")
		}

		startIPInt = (startIPInt << 8) + uint32(startOctet)
		endIPInt = (endIPInt << 8) + uint32(endOctet)
	}

	// Calculate the CIDR mask length
	maskLength := 32 - int(math.Log2(float64(endIPInt-startIPInt+1)))

	if maskLength <= 0 || maskLength > 32 {
		return "", fmt.Errorf("invalid IP address range")
	}

	// Generate the CIDR string
	cidrRange := fmt.Sprintf("%s/%d.%d.%d.%d", startIP, maskLength, startIPInt>>24, (startIPInt>>16)&255, (startIPInt>>8)&255)

	return cidrRange, nil

}

// ExtractIpRange
//
//	@Description: 从文本抽取ip range
//	@return []IpRange
//	@return error
func ExtractIpRange() ([]IpRange, error) {
	//filePath := "tw-hinet.txt"
	filePath := "jp-oracle.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var result []IpRange
	// Iterate over each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into three fields: startIP, endIP, and IPCount
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			fmt.Println("invalid line format:", line)
			continue
		}

		startIP := fields[0]
		endIP := fields[1]
		IPCount := fields[2]
		atoi, err := strconv.Atoi(IPCount)
		if err != nil {
			panic(err)
		}

		// Process the extracted data as needed
		//fmt.Println("Start IP:", startIP)
		//fmt.Println("End IP:", endIP)
		//fmt.Println("IP Count:", IPCount)
		ipRange := IpRange{
			IPStart: startIP,
			IPEnd:   endIP,
			IPCount: atoi,
		}
		result = append(result, ipRange)
	}

	if scanner.Err() != nil {
		fmt.Println("error reading file:", scanner.Err())
		return nil, scanner.Err()
	}
	return result, nil
}

// GetIpListFromIPRange
//
//	@Description: 将iprange 转化为ip列表
//	@param startIP
//	@param endIP
//	@return []string
//	@return error
func GetIpListFromIPRange(startIP, endIP string) ([]string, error) {
	var ipList []string

	// Parse the starting and ending IP addresses
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	// Check if the IP addresses are valid
	if start == nil || end == nil {
		return nil, fmt.Errorf("invalid ip address")
	}

	// Iterate over the IP addresses in the range
	for ip := start; !ip.Equal(end); incIPByOne(ip) {
		ipList = append(ipList, ip.String())
	}
	//add ipend
	//ipList = append(ipList,endIP)

	return ipList, nil
}

// incIPByOne
//
//	@Description: ip递增
//	@param ip
func incIPByOne(ip net.IP) {
	// Increment the IP address by 1
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// 全球ip信息
// http://ipblock.chacuo.net/

// CheckResult
//
//	CheckResult
//	@Description: SNI检查结果
type CheckResult struct {
	Ip        string
	IsProxyIp bool
	Error     error
}

// IpRangeCheckForProxyCDNIP
//
//	@Description: 检查给定的ip范围，并找出反向代理cf的ip
func IpRangeCheckForProxyCDNIP() {
	ipRanges, err := ExtractIpRange()
	if err != nil {
		fmt.Println("Error extracting: ", err)
	}
	for _, value := range ipRanges {
		fmt.Println("process ip range: ", value)
		ipList, err := GetIpListFromIPRange(value.IPStart, value.IPEnd)
		if err != nil {
			fmt.Println("get ip list: ", err)
			break
		}
		batches := divideIntoBatches(ipList, 500)
		for index, v := range batches {
			fmt.Println("正在处理第: ", index+1, "批次")
			CheckIfMatchedCf(v)
		}
	}
}

func main() {
	//IpRangeCheckForProxyCDNIP()

}

// divideIntoBatches
//
//	@Description: 将数据切片分批为指定大小的片段
//	@param data
//	@param batchSize
//	@return [][]string
func divideIntoBatches(data []string, batchSize int) [][]string {
	// 将数据划分成批次
	var batches [][]string
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batches = append(batches, data[i:end])
	}
	return batches
}

// CheckIfMatchedCf
//
//	@Description: 批量检查ip，并判断是否是反向代理cf的ip，最终保存到文本
//	@param sampleIps
func CheckIfMatchedCf(sampleIps []string) {
	resultsChan := make(chan CheckResult, 50)
	waitGroup := &sync.WaitGroup{}

	for _, ip := range sampleIps {
		waitGroup.Add(1)
		go func(ip string) {
			defer waitGroup.Done()
			SNIChecker(ip, "www.cloudflare.com", resultsChan)
		}(ip)
	}
	go func() {
		waitGroup.Wait()
		close(resultsChan)
	}()

	proxyedIps := []string{}
	//channel没关闭 容易阻塞
	for result := range resultsChan {
		if result.Error != nil {
			fmt.Println("checker error:", result.Error)
		}
		if result.IsProxyIp {
			fmt.Printf("ip: %s is a proxy for cloudflare\n", result.Ip)
			proxyedIps = append(proxyedIps, result.Ip)
		}
	}
	fmt.Println("------------------next batch ips------------------")
	var trimedResults []string
	for _, ip := range proxyedIps {
		trimSpace := strings.TrimSpace(ip)
		if trimSpace == "" {
			continue
		}
		trimedResults = append(trimedResults, ip)

	}
	if len(trimedResults) > 0 {
		fmt.Println("------------------write ips to file------------------")

		results := strings.Join(trimedResults, "\n")
		fmt.Println(results)
		err := AppendResultToFile("resultips.txt", results)
		if err != nil {
			fmt.Println("write file error:", err)
		}
	}

}

// AppendResultToFile
//
//	@Description: 文件追加写入
//	@param filePath
//	@param content
//	@return error
func AppendResultToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}

// SNIChecker
//
//	@Description: 使用sni检查代理ip是否反代了指定的sni
//	@param ipStr
//	@param serverName
//	@param resultChan
func SNIChecker(ipStr string, serverName string, resultChan chan CheckResult) {
	dialer := &net.Dialer{
		Timeout: 8 * time.Second,
	}
	// Replace <IP> with the target IP address.
	addr := fmt.Sprintf("%s:443", ipStr)
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		ServerName: serverName,
	})
	if err != nil {
		//fmt.Printf("Error connecting to server: %s to %s\n", err, ipStr)
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

// RemoveDuplicates
//
//	@Description: 移除切片中的重复字符串
//	@param strSlice
//	@return []string
func RemoveDuplicates(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// SimpleSNIChecker
//
//	@Description: 检查是否是sni的代理ip
//	@param ipStr
//	@param sni
//	@return bool
//	@return error
func SimpleSNIChecker(ipStr string, sni string) (bool, error) {
	dialer := &net.Dialer{
		Timeout: 8 * time.Second,
	}
	// Replace <IP> with the target IP address.
	addr := fmt.Sprintf("%s:443", ipStr)
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		ServerName: sni,
	})
	if err != nil {
		fmt.Printf("Error connecting to server: %s to %s\n", err, ipStr)
		return false, err
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
		if strings.Contains(cert.Subject.CommonName, sni) && cert.NotAfter.After(time.Now()) {
			//return true, nil
			isPassed = true
			break
		}

	}
	if isPassed {
		return true, nil
	}
	return false, nil
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
